<template>
  <div
    v-show="isOpen"
    class="shadow-top absolute right-0 bottom-0 left-0 z-10 flex w-full flex-col rounded-t-2xl bg-white dark:border-neutral-950 dark:bg-neutral-800 dark:text-gray-300"
    :style="{
      height: `${currentHeightPx}px`,
      transitionProperty: 'height',
      transitionDuration: isDragging ? '0ms' : '200ms',
      transitionTimingFunction: 'ease-out',
    }"
  >
    <!-- Drag handle bar -->
    <div
      class="relative flex shrink-0 select-none items-center px-2 pt-2 pb-3 touch-none"
      @pointerdown="onDragStart"
    >
      <div class="absolute left-1/2 top-2 h-1.5 w-12 -translate-x-1/2 rounded-full bg-gray-400 dark:bg-gray-600" />
      <div class="flex-1" />
      <button
        type="button"
        class="rounded-full p-1.5 text-gray-500 hover:bg-gray-100 dark:hover:bg-neutral-700"
        :title="$t('close')"
        @click.stop="$emit('close')"
        @pointerdown.stop
      >
        <i-ph-x-bold class="h-4 w-4" />
      </button>
    </div>

    <!-- Scrollable content -->
    <div ref="contentRef" class="min-h-0 flex-1 overflow-y-auto px-4 pb-4">
      <slot />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, onUnmounted, ref } from 'vue';

const props = defineProps<{
  isOpen: boolean;
}>();

defineEmits<{
  close: [];
}>();

const snapPoint = defineModel<'header' | 'half' | 'full'>('snapPoint', {
  default: 'half',
});

const contentRef = ref<HTMLElement>();
const isDragging = ref(false);
const dragCurrentHeight = ref(0);

/** Height in px of each snap point */
function getSnapHeights(): Record<'header' | 'half' | 'full', number> {
  const safeAreaTop =
    parseFloat(getComputedStyle(document.documentElement).getPropertyValue('--safe-area-top')) || 0;
  const appBarSpace =
    parseFloat(getComputedStyle(document.documentElement).getPropertyValue('--app-bar-space')) || 64;
  return {
    header: 96,
    half: Math.round(window.innerHeight * 0.5),
    full: window.innerHeight - safeAreaTop - appBarSpace,
  };
}

const currentHeightPx = computed(() => {
  if (!props.isOpen) return 0;
  if (isDragging.value) return dragCurrentHeight.value;
  return getSnapHeights()[snapPoint.value];
});

// ── Drag state ──────────────────────────────────────────────────────────────
let lastClientY = 0;
let lastTime = 0;
/** Rolling window of px/ms velocities (positive = finger moving down) */
const recentVelocities: number[] = [];
const VELOCITY_THRESHOLD = 0.35; // px/ms

function onDragStart(e: PointerEvent) {
  isDragging.value = true;
  dragCurrentHeight.value = getSnapHeights()[snapPoint.value];
  lastClientY = e.clientY;
  lastTime = Date.now();
  recentVelocities.length = 0;

  window.addEventListener('pointermove', onDragMove, { passive: true });
  window.addEventListener('pointerup', onDragEnd);
  window.addEventListener('pointercancel', onDragEnd);
}

function onDragMove(e: PointerEvent) {
  const now = Date.now();
  const dt = now - lastTime;
  const dy = e.clientY - lastClientY; // positive = finger moves down = sheet shrinks

  if (dt > 0) {
    recentVelocities.push(dy / dt);
    if (recentVelocities.length > 6) recentVelocities.shift();
  }
  lastClientY = e.clientY;
  lastTime = now;

  // Content-scroll priority: when dragging down and content is still scrolled,
  // consume the gesture as a content scroll rather than resizing the sheet.
  if (dy > 0 && contentRef.value && contentRef.value.scrollTop > 0) {
    contentRef.value.scrollTop = Math.max(0, contentRef.value.scrollTop - dy);
    return;
  }

  const { header, full } = getSnapHeights();
  dragCurrentHeight.value = Math.max(header, Math.min(full, dragCurrentHeight.value - dy));
}

function onDragEnd() {
  if (!isDragging.value) return;

  const avgVelocity =
    recentVelocities.length > 0
      ? recentVelocities.reduce((a, b) => a + b, 0) / recentVelocities.length
      : 0;

  const { header, half, full } = getSnapHeights();
  let newSnap: 'header' | 'half' | 'full';

  if (avgVelocity > VELOCITY_THRESHOLD) {
    // Fast swipe down → collapse to header
    newSnap = 'header';
  } else if (avgVelocity < -VELOCITY_THRESHOLD) {
    // Fast swipe up → expand to full
    newSnap = 'full';
  } else {
    // Snap to nearest point
    const h = dragCurrentHeight.value;
    const distances: [number, 'header' | 'half' | 'full'][] = [
      [Math.abs(h - header), 'header'],
      [Math.abs(h - half), 'half'],
      [Math.abs(h - full), 'full'],
    ];
    distances.sort((a, b) => a[0] - b[0]);
    newSnap = distances[0][1];
  }

  snapPoint.value = newSnap;
  isDragging.value = false;
  removeListeners();
}

function removeListeners() {
  window.removeEventListener('pointermove', onDragMove);
  window.removeEventListener('pointerup', onDragEnd);
  window.removeEventListener('pointercancel', onDragEnd);
}

onUnmounted(removeListeners);
</script>
