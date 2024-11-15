

class CategoryRes {
  final int code;
  final CateData data;
  final String msg;

  CategoryRes({required this.code, required this.data, required this.msg});

  factory CategoryRes.fromJson( json) {
    return CategoryRes(
      code: json['code'],
      data: CateData.fromJson(json['data']),
      msg: json['msg'],
    );
  }
}

class CateData {
  final List list;
  final int total;

  CateData({required this.list, required this.total});

  factory CateData.fromJson(Map<String, dynamic> json) {
    var listJson = json['list'] as List;
    List<CategoryItem> itemList = listJson.map((i) => CategoryItem.fromJson(i)).toList();

    return CateData(
      list: itemList,
      total: json['total'],
    );
  }
}


// class CategoryRes {
//   final int code;a
//   final Data data;
//   final List<CategoryItem> list;
//   final int total;

//   Data({required this.list, required this.total});

//   factory Data.fromJson(Map<String, dynamic> json) {
//     var listJson = json['list'] as List;
//     List<CategoryItem> itemList = listJson.map((i) => CategoryItem.fromJson(i)).toList();

//     return Data(
//       list: itemList,
//       total: json['total'],
//     );
//   }
// }

class CategoryItem {
  final int id;
  final String name;
  final int status;
  final int createdAt;
  final int updatedAt;

  CategoryItem({
    required this.id,
    required this.name,
    required this.status,
    required this.createdAt,
    required this.updatedAt,
  });

  factory CategoryItem.fromJson(Map<String, dynamic> json) {
    

    return CategoryItem(
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


