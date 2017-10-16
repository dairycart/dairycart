<template>
    <div>
        <section class="hero is-small">
            <div class="hero-body">
                <div class="container">
                    <h1 class="title">Manage Products</h1>
                </div>
            </div>
        </section>

        <div v-for="productChunk in products" class="columns">
            <div v-for="product in productChunk" v-bind:key="product.sku" class="column">
                <div class="card">
                    <header class="card-header">
                        <a v-bind:href="'product/' + product.sku">
                            <p class="card-header-title">{{product.name}}</p>
                        </a>
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

                            <span class="is-pulled-right">
                                <!-- <span class="icon"> <i class="fa fa-info"></i> </span> -->
                                {{product.quantity}} in stock
                            </span>

                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios';

// const productsURL: string = `//${window.location.hostname}/api/v1/products`;
const productsURL = `//${window.location.hostname}/api/v1/products`;

// function splitIntoSublists(list: Array<Object>): Array<Array<Object>> {
//     let out: Array<Array<Object>> = [];
//     let l: Array<Object> = [];

function splitIntoSublists(list) {
    const out = [];
    let l = [];

    list.forEach((element) => {
        if (l.length === 5) {
            out.push(l);
            l = [];
        }
        l.push(element);
    });
    out.push(l);
    return out;
}

export default {
    data() {
        return {
            loading: true,
            products: [],
            error: null,
        };
    },

    mounted() {
        axios
            .get(productsURL)
            .then((response) => {
                this.products = splitIntoSublists(response.data.data);
                this.loading = false;
            })
            .catch((error) => {
                console.log('error encountered retrieving product data:');
                console.log(error);
                this.loading = false;
                this.error = error;
            });
    },
};
</script>