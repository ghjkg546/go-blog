import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/capchat.dart';
import 'package:flutter_application_2/utils/request.dart';

class WriteCommentPage extends StatefulWidget {
  // const WriteCommentPage({super.key});

  final int id;
  WriteCommentPage({super.key, required this.id});
  @override
  _ComnentPageState createState() => _ComnentPageState();
}

// 评论输入页面
class _ComnentPageState extends State<WriteCommentPage> {
  String captchaId = '';
  String captchaImageUrl = '';

  // 创建 TextEditingController
  final TextEditingController _controller = TextEditingController();
  final TextEditingController captchaController = TextEditingController();
  // 构造函数接收 id 参数
  //  WriteCommentPage({super.key, required this.id});

  void fetchCaptchaImg() async {
    var c1 = await userApi.fetchCaptcha();
    var res = CapchatRes.fromJson(c1).data;
    setState(() {
      captchaId = res.CapchatId;
      captchaImageUrl = Request().baseUrl + res.ImageUrl;
    });
  }

  @override
  void initState() {
    super.initState();
    fetchCaptchaImg();
  }

  void doComment(context) async {
    if (captchaController.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入验证码')),
      );
      return;
    }
    if (_controller.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入评论内容 ')),
      );
      return;
    }

    var res = await userApi.comment(
        widget.id, _controller.text, captchaId, captchaController.text);
    print(res['code']);
    if (res['code'] != 0) {
      fetchCaptchaImg();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(res['msg'])),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('评论成功')),
      );
      Navigator.pop(context, true);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('写评论'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextField(
              controller: _controller, // 绑定 TextEditingController
              maxLines: 5,
              decoration: const InputDecoration(
                hintText: '请输入你的评论...',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 20),
            Row(
              children: [
                if (captchaImageUrl.isNotEmpty)
                  Expanded(
                    flex: 1,
                    child: GestureDetector(
                      onTap: () {
                        fetchCaptchaImg();
                      },
                      child: Image.network(
                        captchaImageUrl,
                        height: 40,
                        fit: BoxFit.cover,
                      ),
                    ),
                  ),
                const SizedBox(width: 8), // 间距
                Expanded(
                  flex: 1,
                  child: TextFormField(
                    controller: captchaController,
                    decoration: const InputDecoration(
                      labelText: '输入验证码',
                      border: OutlineInputBorder(),
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return '请输入验证码';
                      }
                      return null;
                    },
                  ),
                ),
                ElevatedButton(
                  onPressed: () {
                    // 在这里处理评论提交的逻辑
                    doComment(context);
                  },
                  child: const Text('提交评论'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
