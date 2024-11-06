import 'package:flutter/material.dart';
// import 'package:flutter_application_2/detail.dart';
// import 'package:flutter_screenutil/flutter_screenutil.dart';

class MycomponentPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return DefaultTabController(
      length: 6, // 标签页数量
      child: Scaffold(
        appBar: AppBar(
          title: Text("我的组件", style: TextStyle(fontSize: 36)),
          bottom: TabBar(
            indicatorSize: TabBarIndicatorSize.label,
            isScrollable: true,

            //  padding: EdgeInsets.zero,
            indicatorPadding: EdgeInsets.zero,
            //  labelPadding: EdgeInsets.zero,
            indicatorWeight: 4,
            tabs: [
              Tab(child: Text('可爱', style: TextStyle(fontSize: 36))),
              Tab(child: Text('嘲笑', style: TextStyle(fontSize: 36))),
              Tab(child: Text('小号', style: TextStyle(fontSize: 36))),
              Tab(child: Text('中号', style: TextStyle(fontSize: 36))),
            ],
          ),
        ),
        body: TabBarView(
          children: [
            ItemList(),
          ],
        ),
      ),
    );
  }
}

class ItemList extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(24),
      child: ListView.builder(
        itemCount: 10, // 示例数量

        itemBuilder: (context, index) {
          return GestureDetector(
            // onTap: () {
            //   Navigator.push(
            //     context,
            //     MaterialPageRoute(
            //       builder: (context) => DetailPage(index: index),
            //     ),
            //   );
            // },
            child: Row(
              // mainAxisAlignment: MainAxisAlignmentaceAround,
              children: [
                Expanded(child: ItemCard(index: index)),
              ],
            ),
          );
        },
      ),
    );
  }
}

class ItemCard extends StatelessWidget {
  final int index;

  const ItemCard({Key? key, required this.index}) : super(key: key);

  @override
  Widget build(BuildContext context) {
  return Container(
    height: 250,  // 设置高度，Row 会撑满这个高度
    margin: EdgeInsets.all(4),
    padding: EdgeInsets.all(12),
    decoration: BoxDecoration(
      color: Colors.white,
      borderRadius: BorderRadius.circular(10),
      boxShadow: [
        BoxShadow(
          color: Colors.grey.withOpacity(0.5),
          blurRadius: 5,
          spreadRadius: 2,
          offset: Offset(0, 3),
        ),
      ],
    ),
    child: Row(
      crossAxisAlignment: CrossAxisAlignment.stretch, // 让Row的子组件高度撑满
      children: [
        // 左侧的子组件
        Expanded(
          flex: 3,
          child: Container(
            color: Colors.blue,
            child: Center(
              child: Text('左侧'),
            ),
          ),
        ),
        
        // 右侧的子组件
        Expanded(
          flex: 2,
          child: Container(
        
            child:  Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    '绝美',
                    style: TextStyle(fontSize: 26),
                  ),
                  ElevatedButton(
                    onPressed: () {
                      // Your button action here
                    },
                    child: Text('安装',style: TextStyle(fontSize: 26),),
                  ),
                ],
              ),
          ),
        ),

         Expanded(
              flex: 1,
              child: Column(
                 mainAxisAlignment: MainAxisAlignment.end,
                children:[ IconButton(
                  icon: Icon(Icons.delete),
                  onPressed: () {
                    // Delete action here
                  },
                 
                )],
              ),
            ),
      ],
    ),
  );
}


}
