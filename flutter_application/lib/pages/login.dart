import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/register.dart';
import 'package:flutter_application_2/main.dart';
import 'package:flutter_application_2/pages/login.dart';
import 'package:flutter_application_2/pages/my.dart';
import 'package:flutter_application_2/utils/user_preference.dart';

class LoginPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Logintion Page',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: LogintionPage(),
    );
  }
}

class LogintionPage extends StatefulWidget {
  @override
  _LogintionPageState createState() => _LogintionPageState();
}

class _LogintionPageState extends State<LogintionPage> {
  final _formKey = GlobalKey<FormState>();
  String? _username;

  String? _password;

  void _register() {
    if (_formKey.currentState!.validate()) {
      _formKey.currentState!.save();

      dologin(_username, _password);
      // Implement your Logintion logic here
    }
  }

  void dologin(String? _username, String? _password) async {
    var res = await userApi.login(_username, _password);
    if (res['code'] != 0) {
      print(res['message']);
      // var res_msg = res['msg'];
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('登录失败 ' + res['message'])),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('登录成功 $_username')),
      );
      var tokenInfo = UserData.fromJson(res['data']['accessToken']);
      UserPreferences.saveUserInfo(tokenInfo.accessToken, tokenInfo.tokenType);

      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (context) => MyApp(),
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('登录'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                decoration: InputDecoration(labelText: '用户名'),
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
                decoration: InputDecoration(labelText: '密码'),
                obscureText: true,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter your password';
                  }
                  return null;
                },
                onSaved: (value) {
                  _password = value;
                },
              ),
              SizedBox(height: 20),
              ElevatedButton(
                onPressed: _register,
                child: Text('登录'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
