(() => {
    function getNearestChildX(element, x) {
        const elementRect = element.getBoundingClientRect();
        let nearestX = 0;

        for (let child of element.children) {
            const childRect = child.getBoundingClientRect();
            const childX = childRect.left - elementRect.left;

            if (Math.abs(childX - x) < Math.abs(nearestX - x)) {
                nearestX = childX;
            }
        }

        return nearestX;
    }

    function main() {
        const reactRoot = document.getElementById('react-root');
        const carouselElement = document.createElement('div');
        carouselElement.className = 'carousel';
        reactRoot.appendChild(carouselElement);

        const slidesElement = document.createElement('div');
        slidesElement.className = 'slides';
        carouselElement.appendChild(slidesElement);

        for (let i = 0; i < 20; ++i) {
            const slideElement = document.createElement('div');
            slideElement.className = 'slide';
            slidesElement.appendChild(slideElement);
        }

        const estimateElement = document.createElement('div');
        estimateElement.className = 'estimate';
        slidesElement.appendChild(estimateElement);

        const stiffness = 0.0002;
        const dampingRatio = 1;
        const dampingCoefficient = dampingRatio * 2 * Math.sqrt(stiffness);

        let slidesElementOffsetX = 0;
        let slidesElementX = 0;
        let dragVelocity = 0;
        let dragTime = null;
        let simulation = null;

        new Drag(
            slidesElement,
            () => {
                if (simulation !== null) {
                    simulation.stop();
                }

                slidesElementOffsetX = slidesElementX;
            },
            (x, y) => {
                const now = performance.now();

                if (dragTime !== null) {
                    const dx = (slidesElementOffsetX + x) - slidesElementX;
                    const dt = now - dragTime;

                    dragVelocity = dx / dt;
                }

                dragTime = now;

                slidesElementX = slidesElementOffsetX + x;
                slidesElement.style.transform = 'translateX(' + slidesElementX + 'px)';
            },
            () => {
                slidesElementOffsetX = slidesElementX;
                const friction = 0.004;
                const naturalStoppingX = dragVelocity / friction;
                const stoppingXAbsolute = getNearestChildX(slidesElement, -1 * (slidesElementOffsetX + naturalStoppingX));
                const stoppingX = -1 * (slidesElementOffsetX + stoppingXAbsolute);

                const spring = new Spring(stiffness, dampingCoefficient, dragVelocity, -stoppingX);

                simulation = new Simulation(1000 / 60, dt => {
                    spring.step(dt);
                    slidesElementX = slidesElementOffsetX + (spring.displacement + stoppingX);

                    slidesElement.style.transform = 'translateX(' + slidesElementX + 'px)';

                    const stopped = Math.abs(spring.velocity) < 0.001 && Math.abs(spring.acceleration) < 0.001;
                    return !stopped;
                });
            },
        );
    }

    main();
})();
