
class FavRes {
  final int code;
  final FavData data;
  final String msg;

  FavRes({required this.code, required this.data, required this.msg});

  factory FavRes.fromJson( json) {
    return FavRes(
      code: json['code'],
      data: FavData.fromJson(json['data']),
      msg: json['msg'],
    );
  }
}

class FavData {
  final List list;
  final int total;

  FavData({required this.list, required this.total});

  factory FavData.fromJson(Map<String, dynamic> json) {
    var listJson = json['list'] as List;
    List<FavItem> itemList = listJson.map((i) => FavItem.fromJson(i)).toList();

    return FavData(
      list: itemList,
      total: json['total'],
    );
  }
}




class FavItem {
  final int id;
  final String name;
  final int status;
  final int createdAt;
  final int updatedAt;

  FavItem({
    required this.id,
    required this.name,
    required this.status,
    required this.createdAt,
    required this.updatedAt,
  });

  factory FavItem.fromJson(Map<String, dynamic> json) {
    

    return FavItem(
      id: json['id'],
      name: json['name'],
     
      status: json['status'],
   
      createdAt: json['CreatedAt'],
      updatedAt: json['UpdatedAt'],
      // createTimeStr: json['create_time_str'],
      // updateTimeStr: json['update_time_str'],
    );
  }
}


