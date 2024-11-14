import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:get_storage/get_storage.dart';
import 'package:get_x_with_nav/generated/models/user.dart';
import 'package:get_x_with_nav/generated/models/user_model.dart';
import 'package:get_x_with_nav/utils/constants.dart';
import 'package:get_x_with_nav/utils/log.dart';
import 'package:get_x_with_nav/utils/rest_connector.dart';

class LoginAPI {
   login(String username, String password) async {
    final url =
        base_url + "/auth/login";
    log(url);
    try {
      //Response response = await Dio().get(url);
      // Response response = await RestConnector(url: url).po();

      Response res = await RestConnector(
          url: url,
          requestType: "post",
          clearCookies: true,
          data: jsonEncode({
            "username": username,
            "password": password,
          })).getData();
      // print(res);
      //log(response);
        // var tokenInfo = UserRes.fromJson(res);
        print("info");
        print(res.data['data']['accessToken']);
          GetStorage().write("token", res.data['data']['accessToken'].toString());
        return res.data['data']['accessToken'].toString();
      //  final user = User.fromJson(res.data);
      //  return null;
    } catch (e) {
      log(e);
      return null;
    }
  }
}
