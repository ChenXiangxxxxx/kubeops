import {createRouter,createWebHistory} from 'vue-router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import Layout from '@/layout/Layout'

const routes=[
    {
        path: '/',
        redirect: '/home'
    },
    {
        path: "/home",
        component: Layout,
        children: [
            {
                path:"/home",
                name: "概览",
                icon: "fund-outlined",
                meta: {title: "概览"},
                component: ()=> import('@/views/home/Home.vue')
            }
        ]
    },
    {
        path: "/cluster",
        name: "集群",
        component: Layout,
        icon: "cloud-server-outlined",
        children: [
            {
                path: "/cluster/node",
                name: "Node",
                meta: {title: "Node", requireAuth: true},
                component: () => import('@/views/cluster/Node.vue'),
            },
            {
                path: "/cluster/namespace",
                name: "Namespace",
                meta: {title: "Namespace", requireAuth: true},
                component: () => import('@/views/cluster/Namespace.vue'),
            },
            {
                path: "/cluster/pv",
                name: "PV",
                meta: {title: "PV", requireAuth: true},
                component: () => import('@/views/cluster/PV.vue'),
            }
        ]
    },
    {
        path: "/workload",
        name: "工作负载",
        component: Layout,
        icon: "block-outlined",
        children: [
            {
                path: "/workload/pod",
                name: "Pod",
                meta: {title: "Pod", requireAuth: true},
                component: () => import('@/views/workload/Pod.vue'),
            },
            {
                path: "/workload/deployment",
                name: "Deployment",
                meta: {title: "Deployment", requireAuth: true},
                component: () => import('@/views/workload/Deployment.vue'),
            },
            {
                path: "/workload/daemonset",
                name: "DaemonSet",
                meta: {title: "DaemonSet", requireAuth: true},
                component: () => import('@/views/workload/DaemonSet.vue'),
            },
            {
                path: "/workload/statefulset",
                name: "StatefulSet",
                meta: {title: "StatefulSet", requireAuth: true},
                component: () => import('@/views/workload/StatefulSet.vue'),
            },
        ]
    }
]

const router=createRouter({
    history: createWebHistory(),
    routes
})

NProgress.inc(100)
NProgress.configure({ easing:'ease',speed:600,showSpinner:false })
router.beforeEach((to,from,next)=>{
    NProgress.start()

    if (to.meta.title){
        document.title=to.meta.title
    }else{
        document.title="kubeA"
    }
    next()
})
router.afterEach(()=>{
    NProgress.done()
})
export default router

