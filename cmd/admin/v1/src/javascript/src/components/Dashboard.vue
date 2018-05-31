<template>
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
</template>

<script>
export default {
    name: 'HelloWorld',
    data() {
        return {
            msg: 'Welcome to Your Vue.js App',
        };
    },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1,
h2 {
    font-weight: normal;
}

ul {
    list-style-type: none;
    padding: 0;
}

li {
    display: inline-block;
    margin: 0 10px;
}

a {
    color: #42b983;
}
</style>
