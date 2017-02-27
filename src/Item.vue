<template>
    <div>{{item.title}} <button @click="remove">X</button></div>
</template>
<script>
export default {
    name: 'item',
    props: ['item'],
    serverCacheKey: props => props.item.id + '::' + props.item.updated,
    methods: {
        remove () {
            this.$emit('removeClick', this.item.id)
        }
    }
}
</script>