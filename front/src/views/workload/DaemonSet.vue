<template>
  <div>
      <MainHead
          searchDescribe="请输入"
          @searchChange="getSearchValue"
          namespace
          @namespaceChange="getNamespaceValue"
          @dataList="getDaemonSetList"/>
     <a-card :bodyStyle="{padding: '10px'}">
          <a-table
              style="font-size:12px;" 
              :loading="appLoading" 
              :columns="columns" 
              :dataSource="daemonSetList"
              :pagination="pagination"
              @change="handleTableChange">
              <template #bodyCell="{ column, record }">
                  <template v-if="column.dataIndex === 'name'">
                      <span style="font-weight: bold;">{{ record.metadata.name }}</span>
                  </template>
                  <template v-if="column.dataIndex === 'labels'">
                      <div v-for="(val, key) in record.metadata.labels" :key="key">
                          <a-popover>
                              <template #content>
                                  <span>{{ key + ":" +val }}</span>
                              </template>
                              <a-tag style="margin-bottom:5px;cursor:pointer;" color="blue">{{ ellipsis(key + ":" +val, 15) }}</a-tag>
                          </a-popover>
                      </div>
                  </template>
                  <template v-if="column.dataIndex === 'containers'">
                      <span>{{ record.status.numberAvailable>0?record.status.numberAvailable:0  }} / {{ record.status.desiredNumberScheduled>0?record.status.desiredNumberScheduled:0 }} </span>
                  </template>
                  <template v-if="column.dataIndex === 'image'">
                      <div v-for="(val, key) in record.spec.template.spec.containers" :key="key">
                          <a-popover>
                              <template #content>
                                  <span>{{ val.image }}</span>
                              </template>
                              <a-tag style="margin-bottom:5px;cursor:pointer;" color="geekblue">{{ ellipsis(val.image.split('/').pop() ? val.image.split('/').pop() : val.image, 15 ) }}</a-tag>
                          </a-popover>
                      </div>
                  </template>
                  <template v-if="column.dataIndex === 'creationTimestamp'">
                      <a-tag color="gray">{{ timeTrans(record.metadata.creationTimestamp) }}</a-tag>
                  </template>
                  <template v-if="column.key === 'action'">
                      <c-button style="margin-bottom:5px;" class="daemonSet-button" type="primary" icon="form-outlined" @click="getDaemonSetDetail(record)">YML</c-button>
                      <c-button class="daemonSet-button" type="error" icon="delete-outlined" @click="showConfirm('删除', record.metadata.name, delDaemonSet)">删除</c-button>
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
          @ok="updateDaemonSet">
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
import { createVNode, reactive, ref } from 'vue';
import MainHead from '@/components/MainHead';
import httpClient from '@/request';
import common from "@/config";
import { message } from 'ant-design-vue';
import yaml2obj from 'js-yaml';
import json2yaml from 'json2yaml';
import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
import { Modal } from 'ant-design-vue';
export default({
  components: {
      MainHead,
  },
  setup() {
      //表结构
      const columns = ref([
          {
              title: 'DeamonSet名',
              dataIndex: 'name'
          },
          {
              title: '标签',
              dataIndex: 'labels'
          },
          {
              title: '容器组',
              dataIndex: 'containers',
          },
          {
              title: '镜像',
              dataIndex: 'image'
          },
          {
              title: '创建时间',
              dataIndex: 'creationTimestamp'
          },
          {
              title: '操作',
              key: 'action',
              fixed: 'right',
              width: 200
          }
      ])
      //常用项
      const appLoading = ref(false)
      const searchValue = ref('')
      const namespaceValue = ref('')
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
      const daemonSetList = ref([])
      const daemonSetListData = reactive({
          url : common.k8sDaemonSetList,
          params: {
              filter_name: '',
              namespace: '',
              cluster: '',
              page: 1,
              limit: 10
          }
      })
      //详情
      const contentYaml = ref('')
      const yamlModal = ref(false)
      const cmOptions = common.cmOptions
      const daemonSetDetail =  reactive({
          metadata: {}
      })
      const daemonSetDetailData =  reactive({
          url: common.k8sDaemonSetDetail,
          params: {
              daemonset_name: '',
              namespace: '',
              cluster: ''
          }
      })
      //yaml更新
      const updateDaemonSetData = reactive({
          url: common.k8sDaemonSetUpdate,
          params: {
              namespace: '',
              content: '',
              cluster: ''
          }
      })
      //删除
      const delDaemonSetData = reactive({
          url: common.k8sDaemonSetDel,
          params: {
              daemonset_name: '',
              namespace: '',
              cluster: ''
          }
      })

      //json转yaml方法
      function transYaml(content) {
          return json2yaml.stringify(content)
      }
      //yaml转对象
      function transObj(content) {
          return yaml2obj.load(content)
      }
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
          getDaemonSetList()
      }
      function getSearchValue(val) {
          searchValue.value = val
      }
      function getNamespaceValue(val) {
          namespaceValue.value = val
      }
      //编辑器内容变化时触发的方式,用于将更新的内容复制到变量中
      function onChange(val) {
          contentYaml.value = val
      }
      //列表
      function getDaemonSetList() {
          appLoading.value = true
          if (searchValue.value) {
              pagination.currentPage = 1
          }
          daemonSetListData.params.filter_name = searchValue.value
          daemonSetListData.params.namespace = namespaceValue.value
          daemonSetListData.params.cluster = localStorage.getItem('k8s_cluster')
          daemonSetListData.params.page = pagination.currentPage
          daemonSetListData.params.limit = pagination.pagesize
          httpClient.get(daemonSetListData.url, {params: daemonSetListData.params})
          .then(res => {
              //响应成功，获取daemonSet列表和total
              daemonSetList.value = res.data.items
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
      function getDaemonSetDetail(e) {
          appLoading.value = true
          daemonSetDetailData.params.daemonset_name = e.metadata.name
          daemonSetDetailData.params.namespace = namespaceValue.value
          daemonSetDetailData.params.cluster = localStorage.getItem('k8s_cluster')
          httpClient.get(daemonSetDetailData.url, {params: daemonSetDetailData.params})
          .then(res => {
              //daemonSetDetail = Object.assign(daemonSetDetail, res.data)
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
      //更新daemonSet
      function updateDaemonSet() {
          appLoading.value = true
          //将yaml转为json
          let content = JSON.stringify(transObj(contentYaml.value))
          updateDaemonSetData.params.namespace = namespaceValue.value
          updateDaemonSetData.params.content = content
          updateDaemonSetData.params.cluster = localStorage.getItem('k8s_cluster')
          httpClient.put(updateDaemonSetData.url, updateDaemonSetData.params)
          .then(res => {
              message.success(res.msg)
          })
          .catch(res => {
              message.error(res.msg)
          })
          .finally(() => {
              getDaemonSetList()
              yamlModal.value = false
          })
      }
      //删除daemonSet
      function delDaemonSet(name) {
          appLoading.value = true
          delDaemonSetData.params.daemonSet_name = name
          delDaemonSetData.params.namespace = namespaceValue.value
          delDaemonSetData.params.cluster = localStorage.getItem('k8s_cluster')
          httpClient.delete(delDaemonSetData.url, {data: delDaemonSetData.params})
          .then(res => {
              message.success(res.msg)
          })
          .catch(res => {
              message.error(res.msg)
          })
          .finally(() => {
              getDaemonSetList()
          })
      }
      //确认框
      function showConfirm(action, res, fn) {
          Modal.confirm({
              title: '是否继续' + action + "操作? 操作对象：",
              icon: createVNode(ExclamationCircleOutlined),
              content: createVNode('div', {
                  //style: 'color:red;',
              }, res),
              cancelText: '取消',
              okText: '确认',
              onOk() {
                  fn(res)
              },
          })
      }

      return {
          appLoading,
          pagination,
          columns,
          daemonSetList,
          daemonSetDetail,
          contentYaml,
          yamlModal,
          cmOptions,
          timeTrans,
          ellipsis,
          handleTableChange,
          getSearchValue,
          getNamespaceValue,
          getDaemonSetList,
          getDaemonSetDetail,
          onChange,
          updateDaemonSet,
          showConfirm,
          delDaemonSet
      }
  },
})
</script>

<style scoped>
  .daemonSet-button {
      margin-right: 5px;
  }
  .ant-form-item {
      margin-bottom: 20px;
  }
</style>