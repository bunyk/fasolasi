
import matplotlib.pyplot as plt
import numpy as np
import random
t = np.linspace(0, 1, 1000)

# end = 0, 800
# start = 300, 400


for fig in range(5):
    fig = plt.figure()
    ax = fig.add_subplot(1, 1, 1) # nrows, ncols, index

    ax.set_facecolor((0.0, 0.0, 0.0))

    for i in range(30):
        phase = random.random()
        amplitude = random.random()
        frequency = random.random() * 10 - 5
        x = 300 - 100 * t + 100 * np.sin(t * np.pi) * amplitude * np.cos(phase + t * frequency * 3.14)
        y = 400 + 400  * t + 100 * np.sin(t * np.pi) * amplitude * np.sin(phase + t * frequency * 3.14)

        color = (1-amplitude, 1-amplitude, 1-amplitude)
        plt.plot(x, y, c=color)

plt.show()
