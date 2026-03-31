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
      @pointerdown="onHandlePointerDown"
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
import { computed, onMounted, onUnmounted, ref } from 'vue';

const HEADER_SNAP_PX = 96;
const VELOCITY_THRESHOLD = 0.35; // px/ms

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

// ── Snap height calculations ────────────────────────────────────────────────

function getSnapHeights(): Record<'header' | 'half' | 'full', number> {
  // Measure the actual app-bar bottom so we never overlap it
  const appBar = document.getElementById('app-bar');
  const topOffset = appBar ? appBar.getBoundingClientRect().bottom + 8 : 72;

  return {
    header: HEADER_SNAP_PX,
    half: Math.round(window.innerHeight * 0.5),
    full: window.innerHeight - topOffset,
  };
}

const currentHeightPx = computed(() => {
  if (!props.isOpen) return 0;
  if (isDragging.value) return dragCurrentHeight.value;
  return getSnapHeights()[snapPoint.value];
});

// ── Shared drag helpers ─────────────────────────────────────────────────────

let lastClientY = 0;
let lastTime = 0;
const recentVelocities: number[] = [];

function resetDragTracking(clientY: number) {
  lastClientY = clientY;
  lastTime = Date.now();
  recentVelocities.length = 0;
}

/** Record a sample and return the screen-space delta (positive = finger down). */
function trackVelocity(clientY: number): number {
  const now = Date.now();
  const dt = now - lastTime;
  const dy = clientY - lastClientY;

  if (dt > 0) {
    recentVelocities.push(dy / dt);
    if (recentVelocities.length > 6) recentVelocities.shift();
  }

  lastClientY = clientY;
  lastTime = now;
  return dy; // positive = finger moved down
}

/** Snap to the right point based on position + velocity, then exit drag mode. */
function finishDrag() {
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
}

function applyDragDelta(dy: number) {
  const { header, full } = getSnapHeights();
  dragCurrentHeight.value = Math.max(header, Math.min(full, dragCurrentHeight.value - dy));
}

// ── Handle-bar drag (pointer events, touch-none) ────────────────────────────

function onHandlePointerDown(e: PointerEvent) {
  isDragging.value = true;
  dragCurrentHeight.value = getSnapHeights()[snapPoint.value];
  resetDragTracking(e.clientY);

  window.addEventListener('pointermove', onHandlePointerMove, { passive: true });
  window.addEventListener('pointerup', onHandlePointerEnd);
  window.addEventListener('pointercancel', onHandlePointerEnd);
}

function onHandlePointerMove(e: PointerEvent) {
  const dy = trackVelocity(e.clientY);

  // Content-scroll priority: scroll content to top before shrinking
  if (dy > 0 && contentRef.value && contentRef.value.scrollTop > 0) {
    contentRef.value.scrollTop = Math.max(0, contentRef.value.scrollTop - dy);
    return;
  }

  applyDragDelta(dy);
}

function onHandlePointerEnd() {
  finishDrag();
  window.removeEventListener('pointermove', onHandlePointerMove);
  window.removeEventListener('pointerup', onHandlePointerEnd);
  window.removeEventListener('pointercancel', onHandlePointerEnd);
}

// ── Content-area drag (touch events, scroll-first) ──────────────────────────
//
// On the content div the user's finger scrolls content normally.  Only when
// the content reaches a scroll limit (top or bottom) does the gesture switch
// to resizing the sheet.  Once in resize mode the gesture stays there until
// the finger lifts.

let contentResizeActive = false;

function onContentTouchStart(e: TouchEvent) {
  resetDragTracking(e.touches[0].clientY);
  contentResizeActive = false;
}

function onContentTouchMove(e: TouchEvent) {
  const el = contentRef.value;
  if (!el) return;

  const clientY = e.touches[0].clientY;
  const dy = clientY - lastClientY; // positive = finger down

  const atTop = el.scrollTop <= 0;
  const atBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 1;

  // Switch from scroll → resize when a scroll limit is reached
  if (!contentResizeActive) {
    if ((dy > 0 && atTop) || (dy < 0 && atBottom)) {
      contentResizeActive = true;
      isDragging.value = true;
      dragCurrentHeight.value = getSnapHeights()[snapPoint.value];
      el.style.overflowY = 'hidden';
    }
  }

  if (contentResizeActive) {
    e.preventDefault();
    const trackedDy = trackVelocity(clientY);
    applyDragDelta(trackedDy);
  } else {
    // Keep tracking position so velocity is fresh at the moment we switch
    lastClientY = clientY;
    lastTime = Date.now();
  }
}

function onContentTouchEnd() {
  if (contentResizeActive) {
    contentResizeActive = false;
    if (contentRef.value) {
      contentRef.value.style.overflowY = '';
    }
    finishDrag();
  }
}

// ── Lifecycle ───────────────────────────────────────────────────────────────

onMounted(() => {
  const el = contentRef.value;
  if (el) {
    el.addEventListener('touchstart', onContentTouchStart, { passive: true });
    el.addEventListener('touchmove', onContentTouchMove, { passive: false });
    el.addEventListener('touchend', onContentTouchEnd, { passive: true });
    el.addEventListener('touchcancel', onContentTouchEnd, { passive: true });
  }
});

onUnmounted(() => {
  const el = contentRef.value;
  if (el) {
    el.removeEventListener('touchstart', onContentTouchStart);
    el.removeEventListener('touchmove', onContentTouchMove);
    el.removeEventListener('touchend', onContentTouchEnd);
    el.removeEventListener('touchcancel', onContentTouchEnd);
  }
  window.removeEventListener('pointermove', onHandlePointerMove);
  window.removeEventListener('pointerup', onHandlePointerEnd);
  window.removeEventListener('pointercancel', onHandlePointerEnd);
});
</script>
