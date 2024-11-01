import 'package:fluttertoast/fluttertoast.dart';
import 'package:flutter/material.dart';

/// A utility function to show a Toast message.
void showToast({
  required String message,
  Toast length = Toast.LENGTH_SHORT,
  ToastGravity gravity = ToastGravity.BOTTOM,
  Color backgroundColor = Colors.black,
  Color textColor = Colors.white,
  double fontSize = 16.0,
  int timeInSecForIosWeb = 1,
}) {
  Fluttertoast.showToast(
    msg: message,
    toastLength: length,
    gravity: gravity,
    timeInSecForIosWeb: timeInSecForIosWeb,
    backgroundColor: backgroundColor,
    textColor: textColor,
    fontSize: fontSize,
  );
}
