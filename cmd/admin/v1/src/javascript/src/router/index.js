import Vue from 'vue';
import Router from 'vue-router';
import Dashboard from '@/components/Dashboard';
import ProductsPage from '@/components/ProductsPage';

Vue.use(Router);

export default new Router({
    routes: [
        {
            path: '/',
            name: 'Dashboard',
            component: Dashboard,
        },
        {
            path: '/products',
            name: 'Products',
            component: ProductsPage,
        },
    ],
});
