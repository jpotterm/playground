(() => {
    class Simulation {
        constructor(stepInterval, callback) {
            this.stepInterval = stepInterval;
            this.callback = callback;
            this.requestId = null;
            this.previousTime = null;
            this.timeStore = 0;

            this.step = this.step.bind(this);

            this.requestId = requestAnimationFrame(this.step);
        }

        stop() {
            cancelAnimationFrame(this.requestId);
        }

        step(t) {
            if (this.previousTime === null) {
                this.previousTime = t;
                this.requestId = requestAnimationFrame(this.step);
            } else {
                // TODO: maybe have different step and render callbacks so we only render once per raf
                let timeStore = this.timeStore + (t - this.previousTime);
                while (timeStore > this.stepInterval) {
                    timeStore -= this.stepInterval;

                    const shouldContinue = this.callback(this.stepInterval);

                    if (!shouldContinue) return;
                }

                this.timeStore = timeStore;
                this.previousTime = t;
                this.requestId = requestAnimationFrame(this.step);
            }
        }
    }

    class Spring {
        constructor(stiffness, dampingCoefficient, velocity, displacement) {
            this.stiffness = stiffness;
            this.dampingCoefficient = dampingCoefficient;
            this.velocity = velocity;
            this.displacement = displacement;
            this.acceleration = 0;
        }

        step(dt) {
            // We assume a mass of 1 so force = acceleration. You can get the same effect of changing the
            // mass by changing the stiffness and dampingCoefficient
            const acceleration = -1 * this.stiffness * this.displacement - this.dampingCoefficient * this.velocity;

            this.acceleration = acceleration;
            this.velocity = this.velocity + acceleration * dt;
            this.displacement = this.displacement + this.velocity * dt;
        }
    }

    class Drag {
        constructor(element, startCallback, dragCallback, endCallback) {
            this.element = element;
            this.startCallback = startCallback;
            this.dragCallback = dragCallback;
            this.endCallback = endCallback;

            this.offsetX = 0;
            this.offsetY = 0;

            this.onMouseDown = this.onMouseDown.bind(this);
            this.onMouseMove = this.onMouseMove.bind(this);
            this.onMouseUp = this.onMouseUp.bind(this);

            element.addEventListener('mousedown', this.onMouseDown);
        }

        onMouseDown(e) {
            this.startCallback();

            this.offsetX = e.clientX;
            this.offsetY = e.clientY;

            document.addEventListener('mousemove', this.onMouseMove);
            document.addEventListener('mouseup', this.onMouseUp);
        }

        onMouseMove(e) {
            this.dragCallback(e.clientX - this.offsetX, e.clientY - this.offsetY);
        }

        onMouseUp(e) {
            document.removeEventListener('mousemove', this.onMouseMove);
            document.removeEventListener('mouseup', this.onMouseUp);

            this.endCallback();
        }
    }

    function getDampingCoefficient(stiffness, dampingRatio) {
        return dampingRatio * 2 * Math.sqrt(stiffness);
    }

    window.Simulation = Simulation;
    window.Spring = Spring;
    window.Drag = Drag;
    window.getDampingCoefficient = getDampingCoefficient;
})();
