<template>
  <!-- Backdrop overlay -->
  <Transition v-if="showBackdrop" name="backdrop">
    <div
      v-if="isOpen"
      class="fixed inset-0 bg-black bg-opacity-40 z-9 md:hidden"
      @click="
        () => {
          if (closeOnBackdropClick) {
            $emit('close');
          }
        }
      "
    />
  </Transition>

  <!-- Bottom sheet -->
  <div
    ref="sheet"
    class="absolute left-0 right-0 bottom-0 flex flex-col w-full z-100 bg-white shadow-top dark:bg-dark-400 dark:text-gray-300 dark:border-dark-800"
    :class="{
      'rounded-t-2xl': !isFullscreen,
      'rounded-none': isFullscreen,
      'pointer-events-none': !isOpen,
      'p-4': isOpen,
    }"
    :style="sheetStyle"
    @touchstart="handleTouchStart"
    @touchmove="handleTouchMove"
    @touchend="handleTouchEnd"
  >
    <!-- Drag handle -->
    <div
      v-if="!disableResize"
      class="w-full -mt-4 pt-4 pb-4 md:hidden cursor-grab active:cursor-grabbing"
      :class="{ 'pointer-events-auto': isOpen }"
    >
      <div class="flex-shrink-0 bg-gray-500 w-12 h-1.5 rounded-full mx-auto" />
    </div>

    <!-- Content -->
    <div
      ref="content"
      class="flex-1 overflow-y-auto"
      :class="{ 'pointer-events-auto': isOpen }"
      @touchstart="handleContentTouchStart"
      @touchmove="handleContentTouchMove"
    >
      <slot />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref, watch, useTemplateRef } from 'vue';
import { useDomSize } from '~/compositions/useDomSize';

type SnapPoint = number | string;

const props = withDefaults(
  defineProps<{
    isOpen: boolean;
    disableResize?: boolean;
    snapPoints?: SnapPoint[]; // Array of snap points as number (px) or string (px, rem or %)
    currentSnapPoint?: SnapPoint; // Initial snap point value (must be in snapPoints array)
    showBackdrop?: boolean; // Show backdrop overlay
    closeOnBackdropClick?: boolean; // Close when clicking backdrop
    preventClose?: boolean; // Prevent closing by dragging down
  }>(),
  {
    disableResize: false,
    snapPoints: () => ['10%', '50%', '90%'],
    initialSnapPoint: '50%',
    showBackdrop: true,
    closeOnBackdropClick: true,
    preventClose: false,
  },
);

const emit = defineEmits<{
  (event: 'close'): void;
  (event: 'update:current-snap-point', snapPoint: SnapPoint): void;
}>();

const velocityThreshold = 0.5;
const animationDuration = 300;
const animationTimingFunction = 'cubic-bezier(0.32, 0.72, 0, 1)'; // iOS-like spring

const initialSnapPointIndex = computed(() => {
  if (!props.currentSnapPoint) {
    return 0;
  }

  const index = props.snapPoints.indexOf(props.currentSnapPoint);
  if (index === -1) {
    throw new Error('Initial snap point must be one of the defined snap points');
  }

  return index;
});

const sheetRef = useTemplateRef('sheet');
const contentRef = useTemplateRef('content');
const isDragging = ref(false);
const currentHeight = ref(0);
const currentSnapIndex = ref(initialSnapPointIndex.value);
const isAnimating = ref(false);
const isMounted = ref(false);

// Touch tracking
const touchStart = ref({ y: 0, time: 0 });
const lastTouchY = ref(0);
const lastTouchTime = ref(0);

const parentElement = computed(() => sheetRef.value?.parentElement);
const { height: windowHeight } = useDomSize(parentElement);

const remSize = parseFloat(window.getComputedStyle(document.documentElement).fontSize);

// Convert snap point to pixels
function snapPointToPx(snapPoint: SnapPoint): number {
  if (typeof snapPoint === 'string' && snapPoint.endsWith('%')) {
    const percentage = Number.parseFloat(snapPoint) / 100;
    return windowHeight.value * percentage;
  }

  if (typeof snapPoint === 'string' && snapPoint.endsWith('rem')) {
    return Number.parseFloat(snapPoint) * remSize;
  }

  if (typeof snapPoint === 'string' && snapPoint.endsWith('px')) {
    return Number.parseFloat(snapPoint);
  }

  if (typeof snapPoint === 'number') {
    return snapPoint;
  }

  throw new Error('Invalid snap point format');
}

// Get all snap points in pixels
const snapPointsPx = computed(() => {
  return props.snapPoints.map(snapPointToPx).sort((a, b) => a - b);
});

const isFullscreen = computed(() => {
  return currentHeight.value >= windowHeight.value || props.disableResize;
});

// Find nearest snap point based on direction
function findNearestSnapIndex(height: number, velocity: number): number | null {
  const points = snapPointsPx.value;

  // Check velocity for flick gesture
  if (Math.abs(velocity) > velocityThreshold) {
    if (velocity > 0) {
      // Flicking up - go to next higher snap point
      for (let i = 0; i < points.length; i++) {
        if (points[i] > height) {
          return i;
        }
      }
      return points.length - 1;
    } else {
      // Flicking down - go to next lower snap point or close
      for (let i = points.length - 1; i >= 0; i--) {
        if (points[i] < height) {
          return i;
        }
      }
      // Only allow closing if preventClose is false
      return props.preventClose ? 0 : null;
    }
  }

  // No significant velocity - find closest snap point
  let closestIndex = 0;
  let minDistance = Math.abs(height - points[0]);

  for (let i = 1; i < points.length; i++) {
    const distance = Math.abs(height - points[i]);
    if (distance < minDistance) {
      minDistance = distance;
      closestIndex = i;
    }
  }

  // If we're below the minimum snap point by a threshold, close (unless prevented)
  if (height < points[0] * 0.7) {
    return props.preventClose ? 0 : null;
  }

  return closestIndex;
}

// Animate to snap point
function toSnapPoint(snapIndex: number | null, animate: boolean) {
  if (snapIndex === null) {
    currentHeight.value = 0;
    emit('close');
    return;
  }

  if (animate) {
    isAnimating.value = true;
  }

  const height = snapPointsPx.value[snapIndex];
  if (height === undefined) {
    throw new Error('Invalid snap index');
  }
  currentHeight.value = height;
  currentSnapIndex.value = snapIndex;
  const currentSnapPoint = props.snapPoints[snapIndex];
  emit('update:current-snap-point', currentSnapPoint);

  if (animate) {
    setTimeout(() => {
      isAnimating.value = false;
    }, animationDuration);
  }
}
// Touch handlers
function handleTouchStart(e: TouchEvent) {
  if (props.disableResize || !props.isOpen) return;

  // Only handle touches on the drag handle area
  const target = e.target as HTMLElement;
  const isDragHandle = target.closest('.cursor-grab');
  if (!isDragHandle) return;

  isDragging.value = true;
  touchStart.value = {
    y: e.touches[0].clientY,
    time: Date.now(),
  };
  lastTouchY.value = e.touches[0].clientY;
  lastTouchTime.value = Date.now();
}

// Content scroll edge detection
function isScrolledToTop(element: HTMLElement): boolean {
  return element.scrollTop <= 0;
}

function isScrolledToBottom(element: HTMLElement): boolean {
  return Math.abs(element.scrollHeight - element.scrollTop - element.clientHeight) < 1;
}

// Find the scrollable element that was touched
function findScrollableParent(target: HTMLElement): HTMLElement | null {
  let current: HTMLElement | null = target;

  while (current && current !== contentRef.value) {
    const overflowY = window.getComputedStyle(current).overflowY;
    const isScrollable = overflowY === 'auto' || overflowY === 'scroll';
    const hasScroll = current.scrollHeight > current.clientHeight;

    if (isScrollable && hasScroll) {
      return current;
    }

    current = current.parentElement;
  }

  // Check the content container itself
  if (contentRef.value) {
    const overflowY = window.getComputedStyle(contentRef.value).overflowY;
    const isScrollable = overflowY === 'auto' || overflowY === 'scroll';
    const hasScroll = contentRef.value.scrollHeight > contentRef.value.clientHeight;

    if (isScrollable && hasScroll) {
      return contentRef.value;
    }
  }

  return null;
}

let scrollableElement: HTMLElement | null = null;

function handleContentTouchStart(e: TouchEvent) {
  if (props.disableResize || !props.isOpen || !contentRef.value) return;

  // Find the scrollable element at touch point
  scrollableElement = findScrollableParent(e.target as HTMLElement);

  touchStart.value = {
    y: e.touches[0].clientY,
    time: Date.now(),
  };
  lastTouchY.value = e.touches[0].clientY;
  lastTouchTime.value = Date.now();
}

function handleContentTouchMove(e: TouchEvent) {
  if (props.disableResize || !props.isOpen) return;

  const deltaY = e.touches[0].clientY - touchStart.value.y;
  const scrollingDown = deltaY > 0;
  const scrollingUp = deltaY < 0;

  // If there's a scrollable element, check its scroll position
  if (scrollableElement) {
    const shouldDragOnScrollTop = scrollingDown && isScrolledToTop(scrollableElement);
    const shouldDragOnScrollBottom = scrollingUp && isScrolledToBottom(scrollableElement);

    if ((shouldDragOnScrollTop || shouldDragOnScrollBottom) && !isDragging.value) {
      // Start dragging the sheet
      isDragging.value = true;

      // Prevent content scrolling
      e.preventDefault();
    }
  } else if (!isDragging.value) {
    // No scrollable element found, start dragging immediately
    isDragging.value = true;
    e.preventDefault();
  }
}

function handleTouchMove(e: TouchEvent) {
  if (!isDragging.value) return;

  const currentY = e.touches[0].clientY;
  const newHeight = windowHeight.value - currentY;

  // Clamp height with rubber band effect at extremes
  const maxHeight = windowHeight.value;
  const minHeight = 0;

  if (newHeight > maxHeight) {
    // Rubber band at top
    currentHeight.value = maxHeight + (newHeight - maxHeight) * 0.3;
  } else if (newHeight < minHeight) {
    // Rubber band at bottom
    currentHeight.value = newHeight * 0.3;
  } else {
    currentHeight.value = newHeight;
  }

  lastTouchY.value = currentY;
  lastTouchTime.value = Date.now();

  // Prevent scrolling when dragging
  e.preventDefault();
}

function handleTouchEnd() {
  if (!isDragging.value) return;

  isDragging.value = false;

  // Calculate velocity (px/ms)
  const timeDelta = lastTouchTime.value - touchStart.value.time;
  const yDelta = touchStart.value.y - lastTouchY.value; // Positive = dragging up
  const velocity = timeDelta > 0 ? yDelta / timeDelta : 0;

  // Find nearest snap point
  const snapIndex = findNearestSnapIndex(currentHeight.value, velocity);
  toSnapPoint(snapIndex, true);
}

// Computed style for the sheet
const sheetStyle = computed(() => {
  const baseStyle: Record<string, string> = {};

  // Set explicit height
  if (currentHeight.value > 0) {
    baseStyle.height = `${currentHeight.value}px`;
    baseStyle.maxHeight = `${currentHeight.value}px`;
  } else {
    // When closed, set to 0 height
    baseStyle.height = '0px';
    baseStyle.maxHeight = '0px';
    baseStyle.overflow = 'hidden';
  }

  // Add transition for smooth animation
  if ((isAnimating.value || !isDragging.value) && isMounted.value) {
    baseStyle.transition = `height ${animationDuration}ms ${animationTimingFunction}, max-height ${animationDuration}ms ${animationTimingFunction}`;
  }

  // Desktop fullscreen handling
  if (isFullscreen.value) {
    baseStyle.height = '100vh';
    baseStyle.maxHeight = '100vh';
  }

  return baseStyle;
});

// Watch isOpen prop
watch(
  () => props.isOpen,
  (isOpen) => {
    toSnapPoint(isOpen ? initialSnapPointIndex.value : null, true);
  },
);

// Handle resize
function handleResize() {
  if (props.isOpen && !isDragging.value) {
    // Recalculate current snap point height
    currentHeight.value = snapPointsPx.value[currentSnapIndex.value];
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize);

  // Initialize if already open - set height before enabling transitions
  if (props.isOpen) {
    toSnapPoint(initialSnapPointIndex.value, false);
  }

  // Enable transitions after initial setup
  setTimeout(() => {
    isMounted.value = true;
  }, 0);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<style scoped>
.backdrop-enter-active,
.backdrop-leave-active {
  transition: opacity 250ms ease;
}

.backdrop-enter-from,
.backdrop-leave-to {
  opacity: 0;
}
</style>
