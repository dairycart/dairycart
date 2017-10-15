// vue stuff
import Vue from 'vue';
import VueRouter from 'vue-router';

// other dependencies
import axios from 'axios';

const productsURL: string = 'http://localhost:1234/api/v1/products';
const dashboardTemplate: string = `
    <div>
        <section class="hero is-small">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title">Dairycart Dashboard</h1>
                    <h2 class="subtitle">Welcome!</h2>
                </div>
            </div>
        </section>
        <section class="section">
            <div class="columns is-mobile is-multiline">
                <div class="column is-half-desktop is-full-mobile">
                    <section class="panel">
                        <p class="panel-heading">Total Orders</p>
                        <p class="panel-tabs">
                            <a class="is-active" href="#">Past Week</a>
                            <a href="#">Past month</a>
                            <a href="#">Past Quarter</a>
                            <a href="#">Past Year</a>
                            <a href="#">All Time</a>
                        </p>
                        <div class="panel-block">
                            <div id="ordersChart" style="height: 250px;"></div>
                        </div>
                        <div class="panel-block">
                            <button class="button is-default is-outlined is-fullwidth">View Data</button>
                        </div>
                    </section>
                </div>
                <div class="column is-half-desktop is-full-mobile">
                    <section class="panel">
                        <p class="panel-heading">Popular Products</p>
                        <div class="panel-block">
                            <div id="chart2" style="height: 280px;"></div>
                        </div>
                        <div class="panel-block">
                            <button class="button is-default is-outlined is-fullwidth">View Data</button>
                        </div>
                    </section>
                </div>
            </div>
        </section>
    </div>
`;

const productsPageTemplate: string = `
    <div>
        <section class="hero is-small">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title">Manage Products</h1>
                </div>
            </div>
        </section>

        <div v-for="productChunk in products" class="columns">
            <div v-for="product in productChunk" class="column">
                <div class="card">
                    <header class="card-header">
                        <a v-bind:href="'product/' + product.sku"><p class="card-header-title">{{product.name}}</p></a>
                    </header>
                    <div class="card-image">
                        <figure class="image is-4by3">
                            <img v-bind:src="product.imageURL">
                        </figure>
                    </div>
                    <div class="card-content">
                        <div class="panel-block-item">
                            <span class="">
                                <span class="icon">
                                    <i class="fa fa-money"></i>
                                </span>
                                {{product.price}}
                            </span>
                            <!--
                                <span class="is-pulled-right">
                                    <span class="icon"><i class="fa fa-info"></i></span>{{product.quantity}} in stock
                                </span>
                            -->
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
`;

function splitIntoSublists(list: Array<Object>): Array<Array<Object>> {
    let out = [];
    let l = [];
    for (var i = 0; i < list.length; i++) {
        if (l.length === 5) {
            out.push(l);
            l = [];
        }
        l.push(list[i]);
    }
    out.push(l);
    return out;
}

const ProductsPage = {
    template: productsPageTemplate,

    created: function() {
        this.fetchProducts();
    },

    mounted: function() {
        axios
            .get(productsURL)
            .then(response => {
                this.products = splitIntoSublists(response.data.data);
                this.loading = false;
            })
            .catch(function(error) {
                this.loading = false;
                this.error = error;
            });
    },

    data: function() {
        return {
            loading: true,
            products: [],
            error: null,
        };
    },
};

const router = new VueRouter({
    routes: [
        {
            path: '/',
            component: {
                template: dashboardTemplate,
            },
        },
        {
            path: '/product/:sku',
            component: {
                template: dashboardTemplate,
            },
        },
        {
            path: '/products',
            component: ProductsPage,
        },
    ],
});

// 4. Create and mount the root instance.
// Make sure to inject the router with the router option to make the
// whole app router-aware.
const app = new Vue({
    router: router,
    el: '#app',
});
