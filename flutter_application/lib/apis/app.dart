
import 'dart:convert';

import 'package:flutter_application_2/entity/resource_item.dart';

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
  getReasources() async {
    var result = await Request().request(
      "/category/list",
      method: DioMethod.get,
      // data: {"taskuuid": "queryprod", "splist": "66"}
    );

    return result;
  }

  // 获取列表数据
  getInfo() async {
    var result = await Request().request(
      "/user/info",
      method: DioMethod.get,
      // data: {"taskuuid": "queryprod", "splist": "66"}
    );

    return result;
  }

  getListData(int categoryId,int page, String keyword) async {
    var result = await Request().request(
      "/duanju/list?category_id=" + categoryId.toString()+"&page="+page.toString()+"&keyword="+keyword.toString()+"&pageSize=50",
      method: DioMethod.get,
      // data: {"taskuuid": "queryprod", "splist": "66"}
    );
   
    return result;
  }

  getFavList(int categoryId,int page, String keyword) async {
    var result = await Request().request(
      "/user/favlist"+"?page="+page.toString()+"&keyword="+keyword.toString()+"&pageSize=50",
      method: DioMethod.get,
      // data: {"taskuuid": "queryprod", "splist": "66"}
    );
   
    return result;
  }

  // 获取列表数据
  getDetail(int id) async {
    var result = await Request().request("/res/info?id=${id}",
        method: DioMethod.get,);
        // final Map<String, dynamic> jsonMap = jsonDecode(result);
  
// print(res);
  // 输出数据
  // print('Code: ${apiResponse.code}');
  // print('Message: ${apiResponse.msg}');
  // print('Info Name: ${res.data.info.name}');
  // print('Is Favorite: ${res.data.info.isFavorite}');
  // print('Disk Items Array: ${apiResponse.data.info.diskItemsArray.map((item) => item.url).toList()}');
      //  print(result.data['info']['name']);
    return result;
  }

  // 注册
  register(String? username, String? password) async {
    try {
      var result = await Request().request(
        "/auth/register",
        method: DioMethod.post,
        data: {"username": username, "password": password},
      );

      // Assuming `result` is a Map, you can directly access `msg`
      if (result is Map<String, dynamic> && result.containsKey('msg')) {
        print("Message: ${result['msg']}");
        return result;
      } else {
        print("Message not found in the response");
      }
    } catch (e) {
      print("Error during registration: $e");
    }
  }

    // 登录
  login(String? username, String? password) async {
    try {
      var result = await Request().request(
        "/auth/login",
        method: DioMethod.post,
        data: {"username": username, "password": password},
      );
      
      // Assuming `result` is a Map, you can directly access `msg`
      if (result is Map<String, dynamic> && result.containsKey('message')) {
        print("Message: ${result['msg']}");
       
      } 
       return result;
    } catch (e) {
      print("Error during registration: $e");
    }
  }

   // 收藏
  fav(int id) async {
    try {
      var result = await Request().request(
        "/user/fav",
        method: DioMethod.post,
        data: {"id": id},
      );
      
      // Assuming `result` is a Map, you can directly access `msg`
      if (result is Map<String, dynamic> && result.containsKey('message')) {
        print("Message: ${result['msg']}");
       
      } 
       return result;
    } catch (e) {
      print("Error during registration: $e");
    }
  }
}

// 导出全局使用这一个实例
final userApi = UserApi();
