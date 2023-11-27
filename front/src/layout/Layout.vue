<template>
  <a-layout>
        <a-affix>
          <a-layout-header>
            <!--平台信息 float:left居左，并且其他元素也在同一行-->
            <div style="float:left;">
              <img style="height: 40px;margin-bottom: 10px;" :src="kubeLogo"/>
              <span style="font-size: 25px; padding:0 50px 0 20px;font-weight: bold; color: #fff;">kubeA</span>
            </div>
            <!--集群信息-->
            <a-menu
            style="float:left;width:250px;line-height: 64px;"
            v-model:selectdKeys="selectdKeys1"
            theme="dark"
            mode="horizontal">
            <a-menu-item v-for="item in clusterList" :key="item">
              {{ item }}
            </a-menu-item>
            </a-menu>
            <!--用户信息-->
            <div style="float: right;">
              <img style="height: 40px;border-radius: 50%;margin-right: 10px;" :src="avator"/>
              <a-dropdown style="margin-top:50px;" :overLayStyle="{paddintTop:'50px'}">
                <a >
                  Admin
                  <down-outlined/>
                </a>
                <template #overlay>
                  <a-menu>
                    <a-menu-item>
                      <a @click="logout()">退出登录</a>
                    </a-menu-item>
                    <a-menu-item>
                      <a >修改密码</a>
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
          </a-layout-header>
        </a-affix>
        <a-layout style="height: calc(100vh - 68px)">
          <a-layout-sider width="240" v-model:collapsed="collapsed" collapsible>
            <a-menu
            :selectedKeys="selectdKeys2"
            :openKeys="openKeys"
            @openChange="onOpenChange"
            mode="inline"
            :style="{height:'100%',boderRight:0}">
            <template v-for="menu in routers" :key="menu">
              <!--处理无子路由的情况-->
              <a-menu-item 
              v-if="menu.children&&menu.children.length==1"
              :index="menu.children[0].path"
              :key="menu.children[0].path"
              @click="routeChange('item',menu.children[0].path)">
              <template #icon>
                <component :is="menu.children[0].icon"></component>
              </template>
              <span>{{ menu.children[0].name }}</span>
              </a-menu-item>
              <!--处理有子路由的情况-->
              <a-sub-menu 
              v-else-if="menu.children && menu.children.length >1"
              :index="menu.path"
              :key="menu.path">
              <template #icon>
                <component :is="menu.icon"></component>
              </template>
              <template #title>
                <span>
                  <span :class="[collapsed ? 'is-collapse' : '']">{{ menu.name  }}</span>
                </span>
              </template>
              <!--处理子栏目-->
              <a-menu-item
              v-for="child in menu.children"
              :key="child.path"
              :index="child.path"
              @click="routeChange('sub',child.path)">
              <span>{{ child.name }}</span>
            </a-menu-item>
            </a-sub-menu>
            </template>
          </a-menu>
            <!--侧边栏-->
          </a-layout-sider>
          <a-layout style="padding:0 24px">
                <!-- 面包屑 -->
                <a-breadcrumb style="margin: 16px 0">
                    <a-breadcrumb-item>工作台</a-breadcrumb-item>
                    <!-- router.currentRoute.value.matched表示路由的match信息，能拿到父路由和子路由的信息 -->
                    <template v-for="(matched,index) in router.currentRoute.value.matched" :key="index">
                        <a-breadcrumb-item v-if="matched.name">
                            {{ matched.name }}
                        </a-breadcrumb-item>
                    </template>
                </a-breadcrumb>
                <!-- main的部分 -->
            <a-layout-content
            :style="{
              background: 'rgb(31,30,30)',
              padding: '24px',
              margin:0,
              minHeight: '280px',
              overflowY:'auto'}">
              <router-view></router-view>
            </a-layout-content>
            <a-layout-footer style="text-align:center">
              2023 Created by Cx
            </a-layout-footer>
          </a-layout>
        </a-layout>
  </a-layout>
</template>
  
  <script>
  import { ref,onMounted } from 'vue'
  import kubeLogo from '@/assets/k8s-metrics.png'
  import avator from '@/assets/avator.png'
  import { useRouter } from 'vue-router'
  export default({
     setup(){
      const collapsed=ref(false)
      const selectdKeys1=ref([])
      const clusterList = ref([
        'TST-1',
        'TST-2'
      ])
      //侧边栏属性
      const routers = ref([])
      const selectdKeys2 = ref([])
      const openKeys = ref([])
      //通过useRouter方法获取路由配置
      const router = useRouter()
      
      //这里开始是方法
      function logout() {
        //移出用户名
        localStorage.removeItem('username')
        //移出token
        localStorage.removeItem('token')
        //跳转至/login页面
        //router.push('/login')
      }
      //导航栏点击切换页面，以及处理选中的情况
      function routeChange(type,path){
        if (type !='sub'){
          openKeys.value = []
        }
        selectdKeys2.value = [path]
        if (router.currentRoute.value.path !==path){
          router.push(path)
        }
      }
      //用于sub的打开
      function onOpenChange(val){
        const latestOpenKey=val.find(key => openKeys.value.indexOf(key) == -1)
        openKeys.value=latestOpenKey?[latestOpenKey]:[]
      }
      //用于从浏览器地址直接打开后的选中
     
      onMounted(()=>{
        routers.value=router.options.routes
      })
      return{
        collapsed,
        kubeLogo,
        avator,
        selectdKeys1,
        clusterList,
        routers,
        selectdKeys2,
        router,
        openKeys,
        logout,
        routeChange,
        onOpenChange,
      }
     },
  })
  </script>
  

<style scoped>
    .ant-layout-header {
        padding: 0 30px !important;
    }
    .ant-layout-content::-webkit-scrollbar {
        width:6px;
    }
    .ant-layout-content::-webkit-scrollbar-track {
        background-color:rgb(164, 162, 162);
    }
    .ant-layout-content::-webkit-scrollbar-thumb {
        background-color:#666;
    }
    .ant-layout-footer {
        padding: 5px 50px !important;
        color: rgb(239, 239, 239);
    }
    .is-collapse {
        display: none;
    }
    .ant-layout-sider {
        background: #141414 !important;
        overflow-y: auto;
    }
    .ant-layout-sider::-webkit-scrollbar {
        display: none;
    }
    .ant-menu-item {
        margin: 0 !important;
    }
</style>