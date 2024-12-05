import 'dart:convert';

class IndexRes {
  final int code;
  final Data data;
  final String msg;

  IndexRes({
    required this.code,
    required this.data,
    required this.msg,
  });

  factory IndexRes.fromJson(Map<String, dynamic> json) {
    return IndexRes(
      code: json['code'],
      data: Data.fromJson(json['data']),
      msg: json['msg'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'code': code,
      'data': data.toJson(),
      'msg': msg,
    };
  }
}

class Data {
  final List<ListItem> list;
  final int total;

  Data({
    required this.list,
    required this.total,
  });

  factory Data.fromJson(Map<String, dynamic> json) {
    return Data(
      list: (json['list'] as List)
          .map((item) => ListItem.fromJson(item))
          .toList(),
      total: json['total'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'list': list.map((item) => item.toJson()).toList(),
      'total': total,
    };
  }
}

class ListItem {
  final List<ResourceItem> data;
  final String title;

  ListItem({
    required this.data,
    required this.title,
  });

  factory ListItem.fromJson(Map<String, dynamic> json) {
    return ListItem(
      data: (json['data'] as List)
          .map((item) => ResourceItem.fromJson(item))
          .toList(),
      title: json['title'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'data': data.map((item) => item.toJson()).toList(),
      'title': title,
    };
  }
}

class ResourceItem {
  final int id;
  final String name;
  final int categoryId;
  final String description;
  final String coverImg;
  final String diskItems;
  final List<DiskItem> diskItemsArray;
  final String tagIds;
  final String searchId;
  final int status;
  final int isRecommend;
  final int views;
  final int createdAt;
  final int updatedAt;
  final String createTimeStr;
  final String updateTimeStr;
  final String url;
  final bool isFavorite;
  // final Category category;

  ResourceItem({
    required this.id,
    required this.name,
    required this.categoryId,
    required this.description,
    required this.coverImg,
    required this.diskItems,
    required this.diskItemsArray,
    required this.tagIds,
    required this.searchId,
    required this.status,
    required this.isRecommend,
    required this.views,
    required this.createdAt,
    required this.updatedAt,
    required this.createTimeStr,
    required this.updateTimeStr,
    required this.url,
    required this.isFavorite,
    // required this.category,
  });

  factory ResourceItem.fromJson(Map<String, dynamic> json) {
    return ResourceItem(
      id: json['id'],
      name: json['name'],
      categoryId: json['category_id'],
      description: json['description'],
      coverImg: json['cover_img'],
      diskItems: json['disk_items'],
      diskItemsArray: (json['disk_items_array'] as List)
          .map((item) => DiskItem.fromJson(item))
          .toList(),
      tagIds: json['tag_ids'],
      searchId: json['search_id'],
      status: json['status'],
      isRecommend: json['is_recommend'],
      views: json['views'],
      createdAt: json['CreatedAt'],
      updatedAt: json['UpdatedAt'],
      createTimeStr: json['create_time_str'],
      updateTimeStr: json['update_time_str'],
      url: json['url'],
      isFavorite: json['is_favorite'],
      // category: Category.fromJson(json['category']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'category_id': categoryId,
      'description': description,
      'cover_img': coverImg,
      'disk_items': diskItems,
      'disk_items_array': [],
      'tag_ids': tagIds,
      'search_id': searchId,
      'status': status,
      'is_recommend': isRecommend,
      'views': views,
      'CreatedAt': createdAt,
      'UpdatedAt': updatedAt,
      'create_time_str': createTimeStr,
      'update_time_str': updateTimeStr,
      'url': url,
      'is_favorite': isFavorite,
      'category': "",
    };
  }
}

class DiskItem {
  final String url;
  final int type;

  DiskItem({
    required this.url,
    required this.type,
  });

  factory DiskItem.fromJson(Map<String, dynamic> json) {
    return DiskItem(
      url: json['url'],
      type: json['type'],
    );
  }

}