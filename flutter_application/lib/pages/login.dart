import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/register.dart';

import 'package:flutter_application_2/pages/index.dart';
import 'package:flutter_application_2/pages/register.dart';

import 'package:flutter_application_2/utils/user_preference.dart';
import 'package:get/get.dart';



class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  _LogintionPageState createState() => _LogintionPageState();
}

class _LogintionPageState extends State<LoginPage> {
  final _formKey = GlobalKey<FormState>();
  String? _username;

  String? _password;

  void _login() {
    if (_formKey.currentState!.validate()) {
      _formKey.currentState!.save();

      dologin(_username, _password);
      // Implement your Logintion logic here
    }
  }

  void dologin(String? username, String? password) async {
    var res = await userApi.login(username, password);
    if (res['code'] != 0) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('登录失败 ' + res['msg'])),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('登录成功 $username')),
      );
      var tokenInfo = UserData.fromJson(res['data']['accessToken']);

      UserPreferences.saveUserInfo(tokenInfo.accessToken, tokenInfo.tokenType);
      // if (Get.previousRoute.isNotEmpty) {
      //   Get.back();
      // } else {
      Get.to(const IndexPage());
      // 可以在这里处理其他逻辑，比如跳转到首页或退出应用
      //}
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('登录'),
        actions: [
          IconButton(
            icon: const Icon(Icons.favorite_border),
            onPressed: () {
               Navigator.pop(context, 'refresh');
            },
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                decoration: const InputDecoration(labelText: '用户名'),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '请输入用户名';
                  }
                  return null;
                },
                onSaved: (value) {
                  _username = value;
                },
              ),
              TextFormField(
                decoration: const InputDecoration(labelText: '密码'),
                obscureText: true,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '请输入密码';
                  }
                  return null;
                },
                onSaved: (value) {
                  _password = value;
                },
              ),
              const SizedBox(height: 20),
              BrnBigMainButton(
                title: '登录',
                onTap: () {
                  _login();
                },
              ),
              const SizedBox(height: 10,),
              GestureDetector(
                onTap: () => {Get.to(RegisterPage())},
                child: const Text("去注册"))
              
            ],
          ),
        ),
      ),
    );
  }
}
