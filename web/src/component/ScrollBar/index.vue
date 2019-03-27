<template>
    <div class="scroll-container" ref="scrollContainer" @wheel.prevent="handleScroll" >
        <div class="scroll-wrapper" ref="scrollWrapper" :style="{top: top + 'px'}">
            <slot></slot>
        </div>
    </div>
</template>

<script>
const delta = 0;
export default {
    data() {
        return {
            top: 0
        }
    },

    methods: {
        handleScroll(e) {
            const eventDelta = -e.deltaY * 3 || e.wheelDelta;
            const $container = this.$refs.scrollContainer;
            const $containerHeight = $container.offsetHeight;
            const $wrapper = this.$refs.scrollWrapper;
            const $wrapperHeight = $wrapper.offsetHeight;
            if (eventDelta > 0) {
                this.top = Math.min(0, this.top + eventDelta)
            } else {
                if ($containerHeight - delta < $wrapperHeight) {
                    if (this.top < -($wrapperHeight - $containerHeight + delta)) {
                        this.top = this.top
                    } else {
                        this.top = Math.max(this.top + eventDelta, $containerHeight - $wrapperHeight - delta)
                    }
                } else {
                    this.top = 0
                }
            }
        }
    }
}
</script>

<style lang="scss" scoped>
.scroll-container {
    position: relative;
    width: 100%;
    height: 100%;
    overflow: hidden;
    .scroll-wrapper {
        position: absolute;
        width: 100%;
    }
}
</style>
