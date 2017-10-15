"use strict";
exports.__esModule = true;
// vue stuff
var vue_1 = require("vue");
var vue_router_1 = require("vue-router");
// other dependencies
var axios_1 = require("axios");
var productsURL = 'http://localhost:1234/api/v1/products';
var dashboardTemplate = "\n    <div>\n        <section class=\"hero is-small\">\n            <div class=\"hero-body\">\n                <div class=\"container\">\n                    <h1 class=\"title\">Dairycart Dashboard</h1>\n                    <h2 class=\"subtitle\">Welcome!</h2>\n                </div>\n            </div>\n        </section>\n        <section class=\"section\">\n            <div class=\"columns is-mobile is-multiline\">\n                <div class=\"column is-half-desktop is-full-mobile\">\n                    <section class=\"panel\">\n                        <p class=\"panel-heading\">Total Orders</p>\n                        <p class=\"panel-tabs\">\n                            <a class=\"is-active\" href=\"#\">Past Week</a>\n                            <a href=\"#\">Past month</a>\n                            <a href=\"#\">Past Quarter</a>\n                            <a href=\"#\">Past Year</a>\n                            <a href=\"#\">All Time</a>\n                        </p>\n                        <div class=\"panel-block\">\n                            <div id=\"ordersChart\" style=\"height: 250px;\"></div>\n                        </div>\n                        <div class=\"panel-block\">\n                            <button class=\"button is-default is-outlined is-fullwidth\">View Data</button>\n                        </div>\n                    </section>\n                </div>\n                <div class=\"column is-half-desktop is-full-mobile\">\n                    <section class=\"panel\">\n                        <p class=\"panel-heading\">Popular Products</p>\n                        <div class=\"panel-block\">\n                            <div id=\"chart2\" style=\"height: 280px;\"></div>\n                        </div>\n                        <div class=\"panel-block\">\n                            <button class=\"button is-default is-outlined is-fullwidth\">View Data</button>\n                        </div>\n                    </section>\n                </div>\n            </div>\n        </section>\n    </div>\n";
var productsPageTemplate = "\n    <div>\n        <section class=\"hero is-small\">\n            <div class=\"hero-body\">\n                <div class=\"container\">\n                    <h1 class=\"title\">Manage Products</h1>\n                </div>\n            </div>\n        </section>\n\n        <div v-for=\"productChunk in products\" class=\"columns\">\n            <div v-for=\"product in productChunk\" class=\"column\">\n                <div class=\"card\">\n                    <header class=\"card-header\">\n                        <a v-bind:href=\"'product/' + product.sku\"><p class=\"card-header-title\">{{product.name}}</p></a>\n                    </header>\n                    <div class=\"card-image\">\n                        <figure class=\"image is-4by3\">\n                            <img v-bind:src=\"product.imageURL\">\n                        </figure>\n                    </div>\n                    <div class=\"card-content\">\n                        <div class=\"panel-block-item\">\n                            <span class=\"\">\n                                <span class=\"icon\">\n                                    <i class=\"fa fa-money\"></i>\n                                </span>\n                                {{product.price}}\n                            </span>\n                            <!--\n                                <span class=\"is-pulled-right\">\n                                    <span class=\"icon\"><i class=\"fa fa-info\"></i></span>{{product.quantity}} in stock\n                                </span>\n                            -->\n                        </div>\n                    </div>\n                </div>\n            </div>\n        </div>\n    </div>\n";
function splitIntoSublists(list) {
    var out = [];
    var l = [];
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
var ProductsPage = {
    template: productsPageTemplate,
    created: function () {
        this.fetchProducts();
    },
    mounted: function () {
        var _this = this;
        axios_1["default"]
            .get(productsURL)
            .then(function (response) {
            _this.products = splitIntoSublists(response.data.data);
            _this.loading = false;
        })["catch"](function (error) {
            this.loading = false;
            this.error = error;
        });
    },
    data: function () {
        return {
            loading: true,
            products: [],
            error: null
        };
    }
};
var router = new vue_router_1["default"]({
    routes: [
        {
            path: '/',
            component: {
                template: dashboardTemplate
            }
        },
        {
            path: '/product/:sku',
            component: {
                template: dashboardTemplate
            }
        },
        {
            path: '/products',
            component: ProductsPage
        },
    ]
});
// 4. Create and mount the root instance.
// Make sure to inject the router with the router option to make the
// whole app router-aware.
var app = new vue_1["default"]({
    router: router,
    el: '#app'
});
