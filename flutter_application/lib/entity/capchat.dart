


class CapchatRes {
  final int code;
  final CapchatInfoData data;
  final String msg;

  CapchatRes({required this.code, required this.data, required this.msg});

  factory CapchatRes.fromJson( json) {
    return CapchatRes(
      code: json['code'],
      data:  CapchatInfoData.fromJson(json['data']) ,
      msg: json['msg'],
    );
  }
}


class CapchatInfoData {
    String CapchatId;
    
    String ImageUrl;

    CapchatInfoData({
        required this.CapchatId,
        
        required this.ImageUrl,
    });

      // Factory constructor to create Data from JSON
  factory CapchatInfoData.fromJson( json) {
    return CapchatInfoData(
      CapchatId: json['captcha_id'],
      ImageUrl: json['image_url']);
      
  }



}



