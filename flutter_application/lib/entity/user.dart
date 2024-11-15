


class UserRes {
  final int code;
  final UserInfoData? data;
  final String msg;

  UserRes({required this.code, required this.data, required this.msg});

  factory UserRes.fromJson( json) {
    return UserRes(
      code: json['code'],
      data: json['data'] != null ? UserInfoData.fromJson(json['data']) : null,
      msg: json['msg'],
    );
  }
}


class UserInfoData {
    String userId;
    
    String userName;
    int score;
    String avatar;
    UserInfoData({
        required this.userId,
        required this.score,
        required this.avatar,
        
        required this.userName,
    });

      // Factory constructor to create Data from JSON
  factory UserInfoData.fromJson( json) {
    return UserInfoData(
      userId: json['user_id'],
      userName: json['username'],
      avatar: json['avatar'],
      score: json['score']
      );
      
      
  }



}



