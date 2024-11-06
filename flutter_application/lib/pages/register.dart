import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/utils/user_preference.dart';
import 'package:flutter_application_2/entity/register.dart';
import 'package:flutter_application_2/pages/login.dart';

class RegisterPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Registration Page',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: RegistrationPage(),
    );
  }
}

class RegistrationPage extends StatefulWidget {
  @override
  _RegistrationPageState createState() => _RegistrationPageState();
}

class _RegistrationPageState extends State<RegistrationPage> {
  final _formKey = GlobalKey<FormState>();
  String? _username;
  String? _email;
  String? _password;

  void _register() {
    if (_formKey.currentState!.validate()) {
      _formKey.currentState!.save();

      doregister(_username, _password);
      // Implement your registration logic here
    }
  }

  void doregister(String? _username, String? _password) async {
    var res = await userApi.register(_username, _password);
    if (res['code'] != 0) {
      print(res['msg']);
     
      // var res_msg = res['msg'];
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('注册失败 ' + res['msg'])),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('注册成功 $_username')),
      );
       var info  = UserData.fromJson(res['data']);
      print("tokend"+info.accessToken);
        await UserPreferences.saveUserInfo(res, info.accessToken);
   print("User info saved successfully!");
      // Navigator.push(
      //   context,
      //   MaterialPageRoute(
      //     builder: (context) => LoginPage(),
      //   ),
      // );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('注册'),
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
                decoration: InputDecoration(labelText: '邮箱'),
                validator: (value) {
                  // if (value == null || value.isEmpty) {
                  //   return 'Please enter your email';
                  // }
                  // // Simple email validation
                  // if (!RegExp(r'^[^@]+@[^@]+\.[^@]+').hasMatch(value)) {
                  //   return 'Please enter a valid email';
                  // }
                  return null;
                },
                onSaved: (value) {
                  _email = value;
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
                child: Text('注册'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
