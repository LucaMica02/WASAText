<script setup>
import { RouterView } from "vue-router";
</script>
<script>
export default {
  data() {
    return {
      counter: 0,
    };
  },
  methods: {
    incr() {
      this.counter++;
    },
    decr() {
      this.counter--;
    },
    reset() {
      this.counter = 0;
    },
  },
};
</script>

<template>
  <header>APP Header</header>
  <div>
    <button @click="decr()">-</button>
    {{ counter }}
    <button @click="incr()">+</button>
  </div>
  <button @click="reset()"><b>reset</b></button>
  <main><RouterView /></main>
</template>

<style></style>
