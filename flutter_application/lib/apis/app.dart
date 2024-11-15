

import 'dart:io';

import 'package:get/get.dart';

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
      "/duanju/list?category_id=$categoryId&page=$page&keyword=$keyword&pageSize=50",
      method: DioMethod.get,
      // data: {"taskuuid": "queryprod", "splist": "66"}
    );
   
    return result;
  }

  fetchCaptcha() async{
var result = await Request().request(
      "/generate-captcha",
      method: DioMethod.get,
      // data: {"taskuuid": "queryprod", "splist": "66"}
    );
    return result;
  }

  getFavList(int categoryId,int page, String keyword) async {
    var result = await Request().request(
      "/user/favlist?page=$page&keyword=$keyword&pageSize=50",
      method: DioMethod.get,
      // data: {"taskuuid": "queryprod", "splist": "66"}
    );
   
    return result;
  }

  // 获取列表数据
  getDetail(int id) async {
    var result = await Request().request("/res/info?id=$id",
        method: DioMethod.get,);
        // final Map<String, dynamic> jsonMap = jsonDecode(result);
  

    return result;
  }

  // 注册
  register(String? username, String? password,String email, String captchaId, String capchatValue) async {
    try {
      var result = await Request().request(
        "/auth/register",
        method: DioMethod.post,
        data: {"username": username, "password": password, "email": email,"captcha_id":captchaId,"captcha_value":capchatValue},
      );
      return result;
      
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

   // 评论
  comment(int id,String content, String captchaId, String capchatValue) async {
    try {
      var result = await Request().request(
        "/comment/add",
        method: DioMethod.post,
        data: {"resource_item_id": id,"content":content,"captcha_id":captchaId,"captcha_value":capchatValue},
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

  // 签到
  signin() async {
    try {
      var result = await Request().request(
        "/user/sign",
        method: DioMethod.post,
         data: {},
      );
      
     
       return result;
    } catch (e) {
      print("Error during registration: $e");
    }
  }

  getSignStatus() async {
    try {
      var result = await Request().request(
        "/user/signstatus",
        method: DioMethod.get,
        
      );
      
      
       return result;
    } catch (e) {
      print("Error during registration: $e");
    }
  }

  uploadImage(File img ) async {
    /// 开启日志打印
    Request.instance?.openLog();
    try {
       

      var result = await Request().uploadImage(
        img
        // "/image_upload",
        // method: DioMethod.post,
        
      );
      
      
       return result;
    } catch (e) {
      print("Error during registration: $e");
    }
  }


}

// 导出全局使用这一个实例
final userApi = UserApi();
