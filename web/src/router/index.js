import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

const _import = file => () => import('../view/' + file + '.vue')

const appMap = [
    {
        path: '/login',
        name: 'login',
        component: _import('Login'),
    },
]

const routerMap = [
    {
        path: '/',
        component: _import('Layer'),
        meta: {
            single: true,
        },
        redirect: { name: 'dashboard' },
        children: [
            {
                path: 'dashboard',
                name: 'dashboard',
                meta: {
                    title: '控制台',
                    icon: 'icon-dashboard',
                },
                component: _import('Dashboard'),
            },
        ],
    },
]

const router = new Router({
    routes: appMap.concat(routerMap),
    scrollBehavior: () => ({ y: 0 }),
    mode: 'history',
})

export { routerMap }

export default router