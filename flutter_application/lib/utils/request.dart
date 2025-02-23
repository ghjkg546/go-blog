import 'dart:io';

import 'package:dio/dio.dart';
import 'package:flutter_application_2/utils/user_preference.dart';

/// 请求方法:枚举类型
enum DioMethod {
  get,
  post,
  put,
  delete,
  patch,
  head,
}

// 创建请求类：封装dio
class Request {
  /// 单例模式
  static Request? _instance;

  // 工厂函数：执行初始化
  factory Request() => _instance ?? Request._internal();

  // 获取实例对象时，如果有实例对象就返回，没有就初始化
  static Request? get instance => _instance ?? Request._internal();

  /// Dio实例
  static Dio _dio = Dio();

  // String baseUrl = 'https://api.shareziyuan.email/api';
    // String baseUrl = 'http://127.0.0.1:8080/api';
  //  String baseUrl = 'http://192.168.2.236:8080/api';
     String baseUrl = 'http://47.106.155.179:8080/api';



  /// 初始化
  Request._internal() {
    // 初始化基本选项
    BaseOptions options = BaseOptions(
       
             baseUrl: baseUrl,
        connectTimeout: const Duration(seconds: 5),
        receiveTimeout: const Duration(seconds: 5));
    _instance = this;
    // 初始化dio
    _dio = Dio(options);
    // 添加拦截器
    _dio.interceptors.add(InterceptorsWrapper(
        onRequest: _onRequest, onResponse: _onResponse, onError: _onError));
  }

  /// 请求拦截器
  void _onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    // 添加通用参数方法
    // if (!options.path.contains("open")) {
    //   options.queryParameters["userId"] = "xxx";
    // }
     var token = await getToken();

  // Add the Authorization header if the token exists
      if (token != null) {
        options.headers["Authorization"] = "Bearer $token";
      }else{
        options.headers["Authorization"] = "Bearer ";
      }


    // 更多业务需求
    handler.next(options);
  }

  getToken() async {
    var res = await UserPreferences.getAccessToken();
    return res;
  }

  /// 相应拦截器
  void _onResponse(
      Response response, ResponseInterceptorHandler handler) async {
    // 请求成功是对数据做基本处理
    if (response.statusCode == 200) {
      // 处理成功的响应
    } else {
      // 处理异常结果
      print("响应异常: $response");
    }
    handler.next(response);
  }

  /// 错误处理: 网络错误等
  void _onError(DioException error, ErrorInterceptorHandler handler) {
    handler.next(error);
  }

 // 上传图片到服务器
  Future<void> uploadImage(File image) async {
    try {
      String uploadUrl =  baseUrl+'/user/avatar_upload';

      FormData formData = FormData.fromMap({
        'business':'avatar',
        'image': await MultipartFile.fromFile(image.path, filename: 'avatar'),
      });

      var token = await getToken();
      var options = Options(
          headers: {'Content-Type': 'multipart/form-data'},
        );
  // Add the Authorization header if the token exists
      if (token != null) {
        options.headers?["Authorization"] = "Bearer $token";
      }else{
        options.headers?["Authorization"] = "Bearer ";
      }

      Response response = await _dio.post(
        uploadUrl,
        data: formData,
        options: options,
      );

      if (response.statusCode == 200) {
        print('Upload successful: ${response.data}');
      } else {
        print('Upload failed with status: ${response.statusCode}');
      }
    } catch (e) {
      print('Error uploading image: $e');
    }
  }

  /// 请求类：支持异步请求操作
  Future<T> request<T>(
    String path, {
    DioMethod method = DioMethod.get,
    Map<String, dynamic>? params,
    dynamic data,
    CancelToken? cancelToken,
    Options? options,
    ProgressCallback? onSendProgress,
    ProgressCallback? onReceiveProgress,
  }) async {
    const methodValues = {
      DioMethod.get: 'get',
      DioMethod.post: 'post',
      DioMethod.put: 'put',
      DioMethod.delete: 'delete',
      DioMethod.patch: 'patch',
      DioMethod.head: 'head'
    };
    // 默认配置选项
    options ??= Options(method: methodValues[method]);
    try {
      Response response;
      // 开始发送请求
      response = await _dio.request(path,
          data: data,
          queryParameters: params,
          cancelToken: cancelToken,
          options: options,
          onSendProgress: onSendProgress,
          onReceiveProgress: onReceiveProgress);
      return response.data;
    } on DioException catch (e) {
      print("发送请求异常: $e");
      rethrow;
    }
  }

  /// 开启日志打印
  /// 需要打印日志的接口在接口请求前 Request.instance?.openLog();
  void openLog() {
    _dio.interceptors
        .add(LogInterceptor(responseHeader: false, responseBody: true));
  }
}
