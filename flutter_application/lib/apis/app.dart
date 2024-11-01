import 'package:flutter_application_2/entity/category.dart';

import '../utils/request.dart';
 
// 创建一个关于user相关请求的对象
class UserApi {
  /// 单例模式
  static UserApi? _instance;
 
  // 工厂函数：初始化，默认会返回唯一的实例
  factory UserApi() => _instance ?? UserApi._internal();
 
  // 用户Api实例：当访问UserApi的时候，就相当于使用了get方法来获取实例对象，如果_instance存在就返回_instance，不存在就初始化
  static UserApi? get instance => _instance ?? UserApi._internal();
 
  /// 初始化
  UserApi._internal() {
    // 初始化基本选项
  }
 
  /// 获取权限列表
  getUser() async {
    /// 开启日志打印
    Request.instance?.openLog();
 
    /// 发起网络接口请求
    var result = await Request().request('get_user', method: DioMethod.get);
 
    // 返回数据
    return result.data;
  }
 
  // 获取列表数据
  getGoods() async {
    var result = await Request().request("/category/list",
        method: DioMethod.get,
        // data: {"taskuuid": "queryprod", "splist": "66"}
        );
    // 返回数据
    // print("getDetail:$result");
    //     print('ID: ${category.id}, Name: ${category.name}');
    //   });
    return result;
  }

  getListData(int categoryId) async {
    var result = await Request().request("/duanju/list?category_id="+categoryId.toString(),
        method: DioMethod.get,
        // data: {"taskuuid": "queryprod", "splist": "66"}
        );
    // 返回数据
    // print("getDetail:$result");
    
    // print(result);
    // #print(CategoryRes.fromJson(result));
    //  var a = CategoryRes.fromJson(result).data.list;
    // a.forEach((category) {
    //     print('ID: ${category.id}, Name: ${category.name}');
    //   });
    return result;
  }
 
  // 获取列表数据
  getDetail() async {
    var result = await Request().request("/game/gamemgnt",
        method: DioMethod.post,
        data: {"taskuuid": "queryprod", "splist": "66"});
    // 返回数据
    // print("getDetail:$result");
    return result;
  }
}
 
// 导出全局使用这一个实例
final userApi = UserApi();