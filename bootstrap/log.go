package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	dbmodel "github.com/jassue/jassue-gin/app/models"
	"github.com/jassue/jassue-gin/global"
	"github.com/jassue/jassue-gin/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/signal"
	"time"
)

var (
	level   zapcore.Level // zap 日志等级
	options []zap.Option  // zap 配置项
)

func InitializeTele() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}
	fmt.Println(global.App.Config.Search.TeleToken)
	b, err := bot.New(global.App.Config.Search.TeleToken, opts...)
	if nil != err {
		// panics for the sake of simplicity.
		// you should handle this error properly in your code.

		fmt.Println("Failed to create bot:", err)
		return // Exit the goroutine if initialization fails
		//panic(err)
	}

	b.Start(ctx)

}

type Data struct {
	TypeStr string
	Item    dbmodel.ResourceItem
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println(update.Message.Text)

	keyword := update.Message.Text

	db := global.App.DB
	var items []dbmodel.ResourceItem
	query := db.Model(dbmodel.ResourceItem{})
	if keyword != "" {
		query.Where("title LIKE ?", "%"+keyword+"%")
	}

	query.Limit(20).Order("id desc").Find(&items)
	text := ""
	for _, item := range items {
		input := &item
		err1 := json.Unmarshal([]byte(input.DiskItems), &input.DiskItemsArray)
		if err1 != nil {
			fmt.Println("Error decoding JSON:", err1)
			return
		}
		typeStr := ""

		for i := range input.DiskItemsArray {
			menu := input.DiskItemsArray[i]
			typeStr = typeStr + menu.Url + ";"
		}
		text += fmt.Sprintf("标题：%s\n网盘连接：%s\n", item.Title, typeStr)
	}

	if text != "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text,
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "抱歉，暂未收录",
		})
	}
}

func InitializeLog() *zap.Logger {
	// 创建根目录
	createRootDir()

	// 设置日志等级
	setLogLevel()

	if global.App.Config.Log.ShowLine {
		options = append(options, zap.AddCaller())
	}

	// 初始化 zap
	return zap.New(getZapCore(), options...)
}

func createRootDir() {
	if ok, _ := utils.PathExists(global.App.Config.Log.RootDir); !ok {
		_ = os.Mkdir(global.App.Config.Log.RootDir, os.ModePerm)
	}
}

func setLogLevel() {
	switch global.App.Config.Log.Level {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
}

// 扩展 Zap
func getZapCore() zapcore.Core {
	var encoder zapcore.Encoder

	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(global.App.Config.App.Env + "." + l.String())
	}

	// 设置编码器
	if global.App.Config.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, getLogWriter(), level)
}

// 使用 lumberjack 作为日志写入器
func getLogWriter() zapcore.WriteSyncer {
	file := &lumberjack.Logger{
		Filename:   global.App.Config.Log.RootDir + "/" + global.App.Config.Log.Filename,
		MaxSize:    global.App.Config.Log.MaxSize,
		MaxBackups: global.App.Config.Log.MaxBackups,
		MaxAge:     global.App.Config.Log.MaxAge,
		Compress:   global.App.Config.Log.Compress,
	}

	return zapcore.AddSync(file)
}
