
class ApiResponse {
  final int code;
  final String msg;
  final Data data;

  ApiResponse({
    required this.code,
    required this.msg,
    required this.data,
  });

  factory ApiResponse.fromJson( json) {
    return ApiResponse(
      code: json['code'],
      msg: json['msg'],
      data: Data.fromJson(json['data']),
    );
  }
}

class Data {
  final List<dynamic> comments;
  final Info info;

  Data({
    required this.comments,
    required this.info,
  });

  factory Data.fromJson( json) {
    return Data(
      comments: json['comments'] ?? [],
      info: Info.fromJson(json['info']),
    );
  }
}

class Info {
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
  final int views;
  final int createdAt;
  final int updatedAt;
  final String createTimeStr;
  final String updateTimeStr;
  final String url;
  final bool isFavorite;

  Info({
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
    required this.views,
    required this.createdAt,
    required this.updatedAt,
    required this.createTimeStr,
    required this.updateTimeStr,
    required this.url,
    required this.isFavorite,
  });

  factory Info.fromJson(Map<String, dynamic> json) {
    return Info(
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
      views: json['views'],
      createdAt: json['CreatedAt'],
      updatedAt: json['UpdatedAt'],
      createTimeStr: json['create_time_str'],
      updateTimeStr: json['update_time_str'] ?? '',
      url: json['url'] ?? '',
      isFavorite: json['is_favorite'],
    );
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

class CommentEntity {
  final int id;
  final String content;
  final int parentId;
  final int userId;
  final int resourceItemId;
  final int createdAt;
  final int updatedAt;
  final UserEntity user;

  CommentEntity({
    required this.id,
    required this.content,
    required this.parentId,
    required this.userId,
    required this.resourceItemId,
    required this.createdAt,
    required this.updatedAt,
    required this.user,
  });

  // 从 JSON 创建 CommentEntity 对象
  factory CommentEntity.fromJson(Map<String, dynamic> json) {
    return CommentEntity(
      id: json['id'],
      content: json['content'],
      parentId: json['parent_id'],
      userId: json['user_id'],
      resourceItemId: json['resource_item_id'],
      createdAt: json['CreatedAt'],
      updatedAt: json['UpdatedAt'],
      user: UserEntity.fromJson(json['user']),
    );
  }

  // 将 CommentEntity 转换为 JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'content': content,
      'parent_id': parentId,
      'user_id': userId,
      'resource_item_id': resourceItemId,
      'CreatedAt': createdAt,
      'UpdatedAt': updatedAt,
      'user': user.toJson(),
    };
  }
}

class UserEntity {
  final int id;
  final String username;
  final String name;
  final String mobile;
  final String email;
  final String password;
  final int createdAt;
  final int updatedAt;

  UserEntity({
    required this.id,
    required this.username,
    required this.name,
    required this.mobile,
    required this.email,
    required this.password,
    required this.createdAt,
    required this.updatedAt,
  });

  // 从 JSON 创建 UserEntity 对象
  factory UserEntity.fromJson(Map<String, dynamic> json) {
    return UserEntity(
      id: json['id'],
      username: json['username'],
      name: json['name'],
      mobile: json['mobile'],
      email: json['email'],
      password: json['password'],
      createdAt: json['CreatedAt'],
      updatedAt: json['UpdatedAt'],
    );
  }

  // 将 UserEntity 转换为 JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'username': username,
      'name': name,
      'mobile': mobile,
      'email': email,
      'password': password,
      'CreatedAt': createdAt,
      'UpdatedAt': updatedAt,
    };
  }
}
