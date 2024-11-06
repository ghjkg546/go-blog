import 'package:shared_preferences/shared_preferences.dart';

class UserPreferences {

   /// 单例模式
  static UserPreferences? _instance;

  // 工厂函数：初始化，默认会返回唯一的实例
  factory UserPreferences() => _instance ?? UserPreferences._internal();

  // 用户Api实例：当访问UserPreferences的时候，就相当于使用了get方法来获取实例对象，如果_instance存在就返回_instance，不存在就初始化
  static UserPreferences? get instance => _instance ?? UserPreferences._internal();

  /// 初始化
  UserPreferences._internal() {
    // 初始化基本选项
  }
  // Keys for the shared preferences
  static const _keyAccessToken = 'access_token';
  static const _keyUsername = 'username';

  // Save the user's JWT token and username
  static Future<void> saveUserInfo(String accessToken, String username) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_keyAccessToken, accessToken);
    await prefs.setString(_keyUsername, username);
  }

  // Retrieve the access token
   static Future<String?> getAccessToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_keyAccessToken);
  }

  static Future<String?> getAccessTokenDirect() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_keyAccessToken);
  }

  // Retrieve the username
  static Future<String?> getUsername() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_keyUsername);
  }

  // Clear the saved user info (e.g., on logout)
  static Future<void> clearUserInfo() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_keyAccessToken);
    await prefs.remove(_keyUsername);
  }
}
