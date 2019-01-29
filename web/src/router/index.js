import Vue from 'vue'
import Router from 'vue-router'
import i18n from '@/lang'

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
        name: 'main',
        meta: {},
        redirect: { name: 'dashboard' },
        children: [
            {
                path: 'dashboard',
                name: 'dashboard',
                meta: {
                    title: i18n.t('dashboard'),
                    icon: 'icon-dashboard',
                },
                component: _import('Dashboard'),
            },
        ],
    },
    {
        path: '/deploy',
        name: 'deploy',
        component: _import('Layer'),
        meta: {
            title: i18n.t('deploy'),
            icon: 'icon-send',
        },
        children: [
            {
                path: 'apply',
                name: 'deployApply',
                meta: {
                    title: i18n.t('submit_deploy_apply'),
                },
                component: _import('deploy/Apply'),
            },
            {
                path: 'deploy',
                name: 'deployDeploy',
                meta: {
                    title: i18n.t('deploy_manage'),
                },
                component: _import('deploy/Deploy'),
            },
            {
                path: 'release',
                name: 'deployRelease',
                meta: {
                    title: i18n.t('deploying_deploy'),
                    hide: true,
                },
                component: _import('deploy/Release'),
            },
        ],
    },
    {
        path: '/project',
        name: 'project',
        component: _import('Layer'),
        meta: {
            title: i18n.t('project'),
            icon: 'icon-project',
        },
        children: [
            {
                path: 'space',
                name: 'projectSpace',
                meta: {
                    title: i18n.t('space_manage'),
                },
                component: _import('project/Space'),
            },
            {
                path: 'project',
                name: 'projectProject',
                meta: {
                    title: i18n.t('project_manage'),
                    hide: true,
                },
                component: _import('project/Project'),
            },
            {
                path: 'user',
                name: 'projectUser',
                meta: {
                    title: i18n.t('member_manage'),
                    hide: true,
                },
                component: _import('project/Member'),
            },
        ],
    },
    {
        path: '/user',
        name: 'user',
        component: _import('Layer'),
        meta: {
            title: i18n.t('user'),
            icon: 'icon-group',
        },
        children: [
            {
                path: 'group',
                name: 'userGroup',
                meta: {
                    title: i18n.t('role_manage'),
                },
                component: _import('user/Group'),
            },
            {
                path: 'list',
                name: 'userList',
                meta: {
                    title: i18n.t('user_manage'),
                },
                component: _import('user/User'),
            },
        ],
    },
    {
        path: '/server',
        name: 'server',
        component: _import('Layer'),
        meta: {
            title: i18n.t('server'),
            icon: 'icon-server',
        },
        children: [
            {
                path: 'group',
                name: 'serverGroup',
                meta: {
                    title: i18n.t('cluster_manage'),
                },
                component: _import('server/Group'),
            },
            {
                path: 'list',
                name: 'serverList',
                meta: {
                    title: i18n.t('server_manage'),
                },
                component: _import('server/Server'),
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