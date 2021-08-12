import Index from '@/pages/Index'
import Node from '@/pages/Node'
import Rule from '@/pages/Rule'
import System from '@/pages/System'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Index
    },
    {
        path: '/node',
        component: Node,
        name: 'Node'
    },
    {
        path: '/rule',
        component: Rule,
        name: 'Rule'
    },
    {
        path: '/system',
        component: System,
        name: 'System'
    }
]

export default routes