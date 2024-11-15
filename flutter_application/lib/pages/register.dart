import 'package:bruno/bruno.dart';
import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/capchat.dart';

import 'package:flutter_application_2/pages/login.dart';
import 'package:flutter_application_2/utils/request.dart';
import 'package:get/get.dart';

class RegisterPage extends StatefulWidget {
  const RegisterPage({super.key});

  @override
  _RegistrationPageState createState() => _RegistrationPageState();
}

class _RegistrationPageState extends State<RegisterPage> {
  final _formKey = GlobalKey<FormState>();
  String? _username;
  String? _email;
  String? _password;
  String captchaId = '';
  String captchaImageUrl = '';
  final TextEditingController captchaController = TextEditingController();

  void _register() {
    if (_formKey.currentState!.validate()) {
      _formKey.currentState!.save();

      // doregister(_username, _password, _email??"", captchaId, captchaController.text);
      doregister();
      // Implement your registration logic here
    }
  }

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

  void doregister() async {
    var res = await userApi.register(
        _username, _password, _email ?? "", captchaId, captchaController.text);
    if (res['code'] != 0) {
      fetchCaptchaImg();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('注册失败 ' + res['msg'])),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('注册成功 $_username')),
      );

      Get.to(LoginPage());
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('注册'),
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
                    return 'Please enter your username';
                  }
                  return null;
                },
                onSaved: (value) {
                  _username = value;
                },
              ),
              TextFormField(
                decoration: const InputDecoration(labelText: '邮箱'),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '请输入邮箱';
                  }
                  // Simple email validation
                  if (!RegExp(r'^[^@]+@[^@]+\.[^@]+').hasMatch(value)) {
                    return '邮箱格式不正确';
                  }
                  return null;
                },
                onSaved: (value) {
                  _email = value;
                },
              ),
              TextFormField(
                decoration: const InputDecoration(labelText: '密码'),
                obscureText: true,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return '输入密码';
                  }
                  return null;
                },
                onSaved: (value) {
                  _password = value;
                },
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
                ],
              ),
              const SizedBox(height: 16),
              BrnBigMainButton(
                title: '注册',
                onTap: () {
                  _register();
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}
