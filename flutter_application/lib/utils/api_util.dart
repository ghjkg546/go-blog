import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiUtil {
  // 你可以在这里定义基础 URL 和其他通用的设置
  static const String baseUrl = 'https://api.shareziyuan.email/api';

  // Fetch Album 方法
  static Future<Resource> fetchAlbum(int categoryId, int page) async {
    String url = '${baseUrl}list?category_id=$categoryId&pageNum=$page&keyword=';
    
    final response = await http.get(Uri.parse(url));

    if (response.statusCode == 200) {
      return Resource.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
    } else {
      throw Exception('Failed to load album');
    }
  }
}

// 资源模型示例
class Resource {
  // 根据你的实际数据结构定义属性
  Resource.fromJson(Map<String, dynamic> json) {
    // 解析 JSON
  }
}
