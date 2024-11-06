


class RegisterRes {
  final int code;
  final UserData? data;
  final String msg;

  RegisterRes({required this.code, required this.data, required this.msg});

  factory RegisterRes.fromJson( json) {
    return RegisterRes(
      code: json['code'],
      data: json['data'] != null ? UserData.fromJson(json['data']) : null,
      msg: json['msg'],
    );
  }
}


class UserData {
    String accessToken;
    int expiresIn;
    String tokenType;

    UserData({
        required this.accessToken,
        required this.expiresIn,
        required this.tokenType,
    });

      // Factory constructor to create Data from JSON
  factory UserData.fromJson( json) {
    return UserData(
      accessToken: json['access_token'],
      expiresIn: json['expires_in'],
      tokenType: json['token_type'],
    );
  }

  // Method to convert Data to JSON
  // Map<String, dynamic> toJson() {
  //   return {
  //     'access_token': accessToken,
  //     'expires_in': expiresIn,
  //     'token_type': tokenType,
  //   };
  // }

}



