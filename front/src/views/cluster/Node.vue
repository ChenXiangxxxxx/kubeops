<template>
  <div>
      <MainHead
          searchDescribe="请输入"
          @searchChange="getSearchValue"
          @dataList="getNodeList"/>
     <a-card :bodyStyle="{padding: '10px'}">
          <a-table
              style="font-size:12px;" 
              :loading="appLoading" 
              :columns="columns" 
              :dataSource="nodeList"
              :pagination="pagination"
              @change="handleTableChange">
              <template #bodyCell="{ column, record }">
                  <template v-if="column.dataIndex === 'name'">
                      <span style="font-weight: bold;">{{ record.metadata.name }}</span>
                      <br>
                      <span style="color: rgb(84, 138, 238);">{{ record.status.addresses[0].address }}</span>
                  </template>
                  <template v-if="column.dataIndex === 'standard'">
                      <span>{{ record.status.capacity.cpu }}核{{ specTrans(record.status.capacity.memory) }}G</span>
                  </template>
                  <template v-if="column.dataIndex === 'podCidr'">
                      <a-tag color="geekblue">{{ record.spec.podCIDR }}</a-tag>
                  </template>
                  <template v-if="column.dataIndex === 'version'">
                      <span style="color:rgb(13, 173, 231);">{{ record.status.nodeInfo.kubeletVersion }} </span>
                  </template>
                  <template v-if="column.dataIndex === 'creationTimestamp'">
                      <a-tag color="gray">{{ timeTrans(record.metadata.creationTimestamp) }}</a-tag>
                  </template>
                  <template v-if="column.key === 'action'">
                      <c-button style="margin-bottom:5px;" class="node-button" type="primary" icon="form-outlined" @click="getNodeDetail(record)">YML</c-button>
                  </template>
              </template>
          </a-table>
      </a-card>
      <!-- 展示YAML信息的弹框 -->
      <a-modal
          v-model:visible="yamlModal"
          title="YAML信息"
          :confirm-loading="appLoading"
          cancelText="取消"
          okText="更新"
          :ok-button-props="{ disabled: true }">
          <!-- codemirror编辑器 -->
          <!-- border 带边框 -->
          <!-- options  编辑器配置 -->
          <!-- change 编辑器中的内容变化时触发 -->
          <codemirror
              :value="contentYaml"
              border
              :options="cmOptions"
              height="500"
              style="font-size:14px;"
              @change="onChange"
          ></codemirror>
      </a-modal>
  </div>
</template>

<script>
import { reactive, ref } from 'vue';
import MainHead from '@/components/MainHead';
import httpClient from '@/request';
import common from "@/config";
import { message } from 'ant-design-vue';
//import yaml2obj from 'js-yaml';
import json2yaml from 'json2yaml';
export default({
  components: {
      MainHead,
  },
  setup() {
      //表结构
      const columns = ref([
          {
              title: 'Node名',
              dataIndex: 'name'
          },
          {
              title: '规格',
              dataIndex: 'standard'
          },
          {
              title: 'POD-CIDR',
              dataIndex: 'podCidr',
          },
          {
              title: '版本',
              dataIndex: 'version',
          },
          {
              title: '创建时间',
              dataIndex: 'creationTimestamp'
          },
          {
              title: '操作',
              key: 'action',
              fixed: 'right',
              width: 100
          }
      ])
      //常用项
      const appLoading = ref(false)
      const searchValue = ref('')
      //分页
      const pagination = reactive({
          showSizeChanger: true,
          showQuickJumper: true,
          total: 0,
          currentPage: 1,
          pagesize: 10,
          pageSizeOptions: ["10", "20", "50", "100"],
          showTotal: total => `共 ${total} 条`
      })
      //列表
      const nodeList = ref([])
      const nodeListData = reactive({
          url : common.k8sNodeList,
          params: {
              filter_name: '',
              cluster: '',
              page: 1,
              limit: 10
          }
      })
      //详情
      const contentYaml = ref('')
      const yamlModal = ref(false)
      const cmOptions = common.cmOptions
      const nodeDetail =  reactive({
          metadata: {}
      })
      const nodeDetailData =  reactive({
          url: common.k8sNodeDetail,
          params: {
              node_name: '',
              cluster: ''
          }
      })

      //【方法】
      function specTrans(str) {
          if ( str.indexOf('Ki') == -1 ) {
              return str
          }
          let num = str.slice(0,-2) / 1024 / 1024
          return num.toFixed(0)
      }
      //json转yaml方法
      function transYaml(content) {
          return json2yaml.stringify(content)
      }
      //yaml转对象
      // function transObj(content) {
      //     return yaml2obj.load(content)
      // }
      function timeTrans(timestamp) {
          let date = new Date(new Date(timestamp).getTime() + 8 * 3600 * 1000)
          date = date.toJSON();
          date = date.substring(0, 19).replace('T', ' ')
          return date 
      }
      function ellipsis (val, len) {
          return val.length > len ? val.substring(0,len) + '...' : val
      }
      function handleTableChange(val) {
          pagination.currentPage = val.current
          pagination.pagesize = val.pageSize
          getNodeList()
      }
      function getSearchValue(val) {
          searchValue.value = val
      }
      //编辑器内容变化时触发的方式,用于将更新的内容复制到变量中
      function onChange(val) {
          contentYaml.value = val
      }
      //列表
      function getNodeList() {
          appLoading.value = true
          if (searchValue.value) {
              pagination.currentPage = 1
          }
          nodeListData.params.filter_name = searchValue.value
          nodeListData.params.cluster = localStorage.getItem('k8s_cluster')
          nodeListData.params.page = pagination.currentPage
          nodeListData.params.limit = pagination.pagesize
          httpClient.get(nodeListData.url, {params: nodeListData.params})
          .then(res => {
              //响应成功，获取node列表和total
              nodeList.value = res.data.items
              pagination.total = res.data.total
          })
          .catch(res => {
              message.error(res.msg)
          })
          .finally(() => {
              appLoading.value = false
          })
      }
      //详情
      function getNodeDetail(e) {
          appLoading.value = true
          nodeDetailData.params.node_name = e.metadata.name
          nodeDetailData.params.cluster = localStorage.getItem('k8s_cluster')
          httpClient.get(nodeDetailData.url, {params: nodeDetailData.params})
          .then(res => {
              //nodeDetail = Object.assign(nodeDetail, res.data)
              contentYaml.value = transYaml(res.data)
              yamlModal.value = true
          })
          .catch(res => {
              message.error(res.msg)
          })
          .finally(() => {
              appLoading.value = false
          })
      }

      return {
          appLoading,
          pagination,
          columns,
          nodeList,
          nodeDetail,
          contentYaml,
          yamlModal,
          cmOptions,
          timeTrans,
          ellipsis,
          handleTableChange,
          getSearchValue,
          getNodeList,
          getNodeDetail,
          onChange,
          specTrans
      }
  },
})
</script>

<style scoped>
  .node-button {
      margin-right: 5px;
  }
  .ant-form-item {
      margin-bottom: 20px;
      color: rgb(13, 173, 231);
  }
</style>