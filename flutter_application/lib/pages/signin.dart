// ignore_for_file: prefer_const_constructors

import 'package:flutter/material.dart';
import 'package:dio/dio.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/signin.dart';

class SignInPage extends StatefulWidget {
  const SignInPage({super.key});

  @override
  State<SignInPage> createState() => _SignInPageState();
}

class _SignInPageState extends State<SignInPage> {
  final Dio dio = Dio();
  List<int> signStatus = List.filled(7, 0);
  List<int> signScores = List.filled(7, 100); // 签到积分

  @override
  void initState() {
    super.initState();
     fetchSignStatus();
  }

  fetchSignStatus() async {
    var c1 = await userApi.getSignStatus();

var status = SignInResponse.fromJson(c1).data;
    setState(() {
        signStatus = List<int>.from(status);
      });
    
  }

  void doSingin() async {
    var res = await userApi.signin();
    if (res['code'] != 0) {
      
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text( res['msg'])),
      );
    } else {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('签到成功 ')),
      );
      fetchSignStatus();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('周签到'),
      ),
      body: Column(
        children: [
          const SizedBox(height: 20),
          const SizedBox(height: 20),
          Wrap(
            spacing: 10, // 水平方向的间距
            runSpacing: 10, // 垂直方向的间距
            alignment: WrapAlignment.center,
            children: List.generate(7, (index) {
              return Column(
                children: [
                  Container(
                    margin: const EdgeInsets.all(4),
                    width: 50,
                    height: 50,
                    decoration: BoxDecoration(
                      color:
                          signStatus[index] == 1 ? Colors.green : Colors.grey,
                      shape: BoxShape.circle,
                    ),
                    child: Center(child: Text(index==0?'周日':'周${index }')),
                  ),
                  const SizedBox(height: 5),
                  Text(
                    '${signScores[index]}',
                    style: TextStyle(
                      fontSize: 14,
                      color:
                          signStatus[index] == 1 ? Colors.green : Colors.black,
                    ),
                  ),
                ],
              );
            }),
          ),
          SizedBox(
            height: 20,
          ),
          ElevatedButton(
            onPressed: doSingin,
            child: const Text('签到'),
          ),
        ],
      ),
    );
  }
}
