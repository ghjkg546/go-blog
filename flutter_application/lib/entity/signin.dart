class SignInResponse {
  final int code;
  final List<int> data;
  final String msg;

  SignInResponse({
    required this.code,
    required this.data,
    required this.msg,
  });

  // 从 JSON 数据解析
  factory SignInResponse.fromJson(Map<String, dynamic> json) {
    return SignInResponse(
      code: json['code'] ?? 0,
      data: List<int>.from(json['data'] ?? []),
      msg: json['msg'] ?? '',
    );
  }

  // 转换为 JSON 数据
  Map<String, dynamic> toJson() {
    return {
      'code': code,
      'data': data,
      'msg': msg,
    };
  }
}
