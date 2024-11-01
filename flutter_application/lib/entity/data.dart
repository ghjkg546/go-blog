import 'package:flutter_application_2/entity/category.dart';


class DataRes {
  final int code;
  final ItemData data;
  final String msg;

  DataRes({required this.code, required this.data, required this.msg});

  factory DataRes.fromJson( json) {
    return DataRes(
      code: json['code'],
      data: ItemData.fromJson(json['data']),
      msg: json['msg'],
    );
  }
}

// class Data {
//   final List<Item> list;
//   final int total;

//   Data({required this.list, required this.total});

//   factory Data.fromJson(Map<String, dynamic> json) {
//     var listJson = json['list'] as List;
//     List<Item> itemList = listJson.map((i) => Item.fromJson(i)).toList();

//     return Data(
//       list: itemList,
//       total: json['total'],
//     );
//   }
// }

class ItemData {
  final List list;
  final int total;

  ItemData({required this.list, required this.total});

  factory ItemData.fromJson(Map<String, dynamic> json) {
    var listJson = json['list'] as List;
    List<Item> itemList = listJson.map((i) => Item.fromJson(i)).toList();

    return ItemData(
      list: itemList,
      total: json['total'],
    );
  }
}

class Item {
  final int id;
  final String name;
  final int categoryId;
  final String description;
  final String coverImg;
  final String diskItems;
  final List<DiskItem> diskItemsArray;
  final String tagIds;
  final int status;
  final int views;
  final int createdAt;
  final int updatedAt;
  final String createTimeStr;
  final String updateTimeStr;

  Item({
    required this.id,
    required this.name,
    required this.categoryId,
    required this.description,
    required this.coverImg,
    required this.diskItems,
    required this.diskItemsArray,
    required this.tagIds,
    required this.status,
    required this.views,
    required this.createdAt,
    required this.updatedAt,
    required this.createTimeStr,
    required this.updateTimeStr,
  });

  factory Item.fromJson(Map<String, dynamic> json) {
    var diskItemsArrayJson = json['disk_items_array'] as List;
    List<DiskItem> diskItemsList =
        diskItemsArrayJson.map((i) => DiskItem.fromJson(i)).toList();

    return Item(
      id: json['id'],
      name: json['name'],
      categoryId: json['category_id'],
      description: json['description'],
      coverImg: json['cover_img'],
      diskItems: json['disk_items'],
      diskItemsArray: diskItemsList,
      tagIds: json['tag_ids'],
      status: json['status'],
      views: json['views'],
      createdAt: json['CreatedAt'],
      updatedAt: json['UpdatedAt'],
      createTimeStr: json['create_time_str'],
      updateTimeStr: json['update_time_str'],
    );
  }
}

class DiskItem {
  final String url;
  final int type;

  DiskItem({required this.url, required this.type});

  factory DiskItem.fromJson(Map<String, dynamic> json) {
    return DiskItem(
      url: json['url'],
      type: json['type'],
    );
  }
}
