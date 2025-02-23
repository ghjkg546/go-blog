import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import 'package:get_x_with_nav/routes/app_pages.dart';
import 'package:get_x_with_nav/utils/colors.dart';
import 'package:get_x_with_nav/utils/constants.dart';
import 'package:get_x_with_nav/utils/images.dart';
import 'package:get_x_with_nav/utils/widgets.dart';

Widget getDrawerItem(String icon, String name, int pos) {
  return GestureDetector(
    onTap: () {
      if (pos == 5) {
        print("logout clicked....");
        GetStorage().write("myUser", "");
        Get.toNamed(Routes.LOGIN);
      }
    },
    child: Container(
      //color: selectedItem == pos ? t2_colorPrimaryLight : t2_white,
      //color: white,

      padding: EdgeInsets.fromLTRB(10, 6, 10, 6),
      child: Row(
        children: <Widget>[
          //SvgPicture.asset(icon, width: 20, height: 20),
          SizedBox(width: 10),
          text(name,
              textColor: colorPrimary,
              //selectedItem == pos ? t2_colorPrimary : t2TextColorPrimary,
              fontSize: textSizeLargeMedium,
              fontFamily: fontMedium)
        ],
      ),
    ),
  );
}

Widget leftNavi(String name) {
  return SizedBox(
    width: Get.width * 0.85,
    height: Get.height,
    child: Drawer(
      elevation: 8,
      child: SingleChildScrollView(
        child: Container(
          width: Get.width,
          //color: white,
          child: Column(
            children: <Widget>[
              Padding(
                  padding: const EdgeInsets.only(top: 40, right: 20),
                  child: Container(
                    padding: EdgeInsets.fromLTRB(20, 20, 20, 20),
                    decoration: new BoxDecoration(
                        color: colorPrimary,
                        borderRadius: new BorderRadius.only(
                            bottomRight: const Radius.circular(24.0),
                            topRight: const Radius.circular(24.0))),
                    /*User Profile*/
                    child: Row(
                      children: <Widget>[
                        CircleAvatar(
                            backgroundImage:
                                Image.asset(t2_profile, height: 100, width: 100)
                                    .image,
                            radius: 40),
                        SizedBox(width: 16),
                        Expanded(
                          child: Container(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              mainAxisAlignment: MainAxisAlignment.center,
                              children: <Widget>[
                                text("left_welcome".tr + name,
                                    textColor: white,
                                    fontFamily: fontBold,
                                    fontSize: textSizeNormal),
                                SizedBox(height: 8),
                                text("user_email@name.com",
                                    textColor: white, fontSize: textSizeMedium),
                              ],
                            ),
                          ),
                        )
                      ],
                    ),
                  )),
              SizedBox(height: 10),
              getDrawerItem("t2_user", "left_profile".tr, 1),
              getDrawerItem("t2_chat", "left_noti".tr, 2),
              getDrawerItem("t2_report", "left_reports".tr, 3),
              getDrawerItem("t2_settings", "left_settings".tr, 4),
              getDrawerItem("t2_logout", "left_signout".tr, 5),
              SizedBox(height: 10),
              Divider(color: view_color, height: 1),
              SizedBox(height: 10),
              getDrawerItem("t2_share", "left_share".tr, 6),
              getDrawerItem("t2_help", "left_help".tr, 7),
              SizedBox(height: 10),
            ],
          ),
        ),
      ),
    ),
  );
}
