import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_application_2/apis/app.dart';
import 'package:flutter_application_2/entity/category.dart';
import 'package:flutter_application_2/entity/data.dart';
import 'package:flutter_application_2/pages/detail.dart';
import 'package:flutter_application_2/item_card.dart';
import 'package:http/http.dart' as http;
import 'package:number_paginator/number_paginator.dart';

// Future<Resource> fetchAlbum(int categoryId, int page) async {
//   String url = 'https://api.shareziyuan.email/api/res/list?category_id=' +
//       categoryId.toString() +
//       '&pageNum=' +
//       page.toString() +
//       '&keyword=';
//   final response = await http.get(Uri.parse(url));
//   if (response.statusCode == 200) {
//     return Resource.fromJson(jsonDecode(response.body) as Map<String, dynamic>);
//   } else {
//     // If the server did not return a 200 OK response,
//     // then throw an exception.
//     throw Exception('Failed to load album');
//   }
// }

// fetchCategory() async {
//   print("请求分类");
//   var res = await userApi.getGoods();

//   return res;

//   // print(response);
//   // String url = 'https://api.shareziyuan.email/api/category/list';
//   // print(url);
//   // final response = await http.get(Uri.parse(url));
//   // print(CategoryRes.fromJson(
//   //       jsonDecode(response) as Map<String, dynamic>));
//   // return CategoryRes.fromJson(
//   //       jsonDecode(response) as Map<String, dynamic>);
//   // if (response.statusCode == 200) {
//   //   return CategoryRes.fromJson(
//   //       jsonDecode(response.body) as Map<String, dynamic>);
//   // } else {
//   //   throw Exception('Failed to load category');
//   // }
// }



// void main() => runApp(const MyApp());

class ListPage extends StatefulWidget {
   const ListPage({super.key});

  // @override
   State<ListPage> createState() => _MyAppState();
}

class _MyAppState extends State<ListPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  // late Future<Resource> futureReasorce;
  // late Future<CategoryRes> futureCategory;
  // final List<String> _categories = ['pc游戏', '安卓游戏'];
  List cates = [];
  List resource_items  =[];
  // int _page = 1;
  // int _categoryId = 0;
  @override
  void initState() {
    super.initState();

    // futureReasorce = fetchAlbum(_categoryId, _page);
    getGoods();
  }

  // 获取列表
  getGoods() async {
    var c1 = await userApi.getGoods();
    
   
    var cates1 = CategoryRes.fromJson(c1).data.list;
  

    setState(() {
      cates = cates1;
      _tabController = TabController(length: cates.length, vsync: this);
      if(cates.length<=0){
        return;
      }
      getListData(cates[0].id);
      // print("cates")
      // print(fu)
      _tabController.addListener(() {
        if (_tabController.indexIsChanging) {
          // Fetch new data when tab changes
          setState(() {
           
            getListData(cates[_tabController.index].id);
            // print(object)
            // futureReasorce = fetchAlbum(_categoryId, _page);
          });
        }
      });
      
    });
  }

  getListData(int categoryId) async {
    var c1 = await userApi.getListData(categoryId);
   
   
      var cates1 = DataRes.fromJson(c1).data.list;
      // resource_items = cates1;

      setState(() {
        resource_items = cates1;
      });
    // print(cates1);

   
  }

   @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: '资源列表',
        theme: ThemeData(
          colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        ),
        home: Scaffold(
            resizeToAvoidBottomInset: false,
            appBar: AppBar(
              
              bottom: TabBar(
                controller: _tabController,
                tabs:
                    cates.map((category) => Tab(text: category.name)).toList(),
              ),
            ),
            body: Column(
              children: [
                Expanded(
                  child: TabBarView(
                    controller: _tabController,
                    children: resource_items.map((_item) {
                // 使用 GridView.builder 替换内容
                return GridView.builder(
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 2, // 每行显示 2 个项
                    childAspectRatio: 1, // 设置子项的宽高比
                  ),
                  itemCount: resource_items.length, // 这里你可以根据具体数据调整数量
                  itemBuilder: (context, index) {
                    return ItemCard(item: resource_items[index]); // 传入每个项
                  },
                );
              }).toList(),
                // 每个 Tab 对应的内容
                 //       return GridView.builder(
                    //         padding: const EdgeInsets.all(8.0),
                    //         gridDelegate:
                    //             const SliverGridDelegateWithFixedCrossAxisCount(
                    //           crossAxisCount: 2,
                    //           crossAxisSpacing: 8.0,
                    //           mainAxisSpacing: 8.0,
                    //         ),
                    //         itemCount: resource.data.list.length,
                    //         itemBuilder: (context, index) {
                    //           final item = resource.data.list[index];
                    
                    // }).toList(),
                // return Center(child: Text('内容 for ${item.name}'));
              // }).toList(),
                    
                    
                    // _categories.map((category) {
                    //   return FutureBuilder<Resource>(
                    //     future: futureReasorce,
                    //     builder: (context, snapshot) {
                    //       if (snapshot.connectionState ==
                    //           ConnectionState.waiting) {
                    //         return const Center(
                    //             child: CircularProgressIndicator());
                    //       } else if (snapshot.hasError) {
                    //         return Center(
                    //             child: Text('Error: ${snapshot.error}'));
                    //       } else if (!snapshot.hasData ||
                    //           snapshot.data!.data.list.isEmpty) {
                    //         return const Center(
                    //             child: Text('No data available'));
                    //       }

                    //       final resource = snapshot.data!;

                    //       return GridView.builder(
                    //         padding: const EdgeInsets.all(8.0),
                    //         gridDelegate:
                    //             const SliverGridDelegateWithFixedCrossAxisCount(
                    //           crossAxisCount: 2,
                    //           crossAxisSpacing: 8.0,
                    //           mainAxisSpacing: 8.0,
                    //         ),
                    //         itemCount: resource.data.list.length,
                    //         itemBuilder: (context, index) {
                    //           final item = resource.data.list[index];
                    //           return Card(
                    //             elevation: 4.0,
                    //             child: Column(
                    //               crossAxisAlignment: CrossAxisAlignment.start,
                    //               children: [
                    //                 GestureDetector(
                    //                   onTap: () {
                    //                     Navigator.push(
                    //                       context,
                    //                       MaterialPageRoute(
                    //                         builder: (context) =>
                    //                             SecondRoute(item: item),
                    //                       ),
                    //                     );
                    //                   },
                    //                   child: SizedBox(
                    //                     height:
                    //                         150, // Specify a height for the image
                    //                     width: double.infinity,
                    //                     child: item.coverImg.isNotEmpty
                    //                         ? Image.network(
                    //                             item.coverImg,
                    //                             fit: BoxFit.cover,
                    //                             loadingBuilder: (context, child,
                    //                                 loadingProgress) {
                    //                               if (loadingProgress == null) {
                    //                                 return child;
                    //                               } else {
                    //                                 return const Center(
                    //                                     child:
                    //                                         CircularProgressIndicator());
                    //                               }
                    //                             },
                    //                             errorBuilder: (context, error,
                    //                                 stackTrace) {
                    //                               return Image.asset(
                    //                                   'images/noimage.png');
                    //                             },
                    //                           )
                    //                         : Image.asset('images/noimage.png',
                    //                             fit: BoxFit.cover),
                    //                   ),
                    //                 ),
                    //                 Padding(
                    //                   padding: const EdgeInsets.all(8.0),
                    //                   child: Text(
                    //                     item.name,
                    //                     style: const TextStyle(
                    //                       fontWeight: FontWeight.bold,
                    //                       fontSize: 16,
                    //                     ),
                    //                   ),
                    //                 ),
                    //               ],
                    //             ),
                    //           );
                    //         },
                    //       );
                    //     },
                    //   );
                    // }).toList(),
                  ),
                ),
                // NumberPaginator(
                //   numberPages: 10,
                //   onPageChange: (int index) {
                //     setState(() {
                //       futureReasorce = fetchAlbum(_categoryId, index + 1);
                //     });
                //   },
                // ),
              ],
            )
            
            )
            );
  }
}
