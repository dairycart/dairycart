<template>
    <div>
        <section class="hero is-small">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title">Manage Products</h1>
                </div>
            </div>
        </section>

        <div v-for="product in products" :key="product.SKU" class="columns">
            <div class="column">
                <div class="card">
                    <header class="card-header">
                        <a v-bind:href="'/product/' + product.SKU">
                            <p class="card-header-title">{{product.name}}</p>
                        </a>
                    </header>
                    <div class="card-image">
                        <figure class="image is-4by3">
                            <img v-bind:src="imageURL">
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

                            <span class="is-pulled-right">
                                <span class="icon">
                                    <i class="fa fa-info"></i>
                                </span>{{product.quantity}} in stock
                            </span>

                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'Products',

    template: productsPageTemplate,

    data() {
        return {
            loading: true,
            products: [],
            error: null,
        };
    },

    created: function() {
        this.fetchProducts();
    },

    methods: {
        fetchProducts: function() {
            const productURL =
                'http://admin.dairycart.com/api/v1/products';

            axios
                .get(productURL)
                .then(function(response) {
                    products = response.data.data;
                    this.products = splitIntoSublists(products);
                    this.loading = false;
                })
                .catch(function(error) {
                    this.error = error;
                });
        },
    },
    created() {
        this.fetchData()
    },
}
</script>