### Real-World System Failure: The 2024 CrowdStrike Outage

* **The Update:** CrowdStrike released a routine threat-detection configuration file update to its systems.
* **The Fault (Software Bug):** A small bug occurred when the system tried to open the file, causing a "memory read out of bounds" error. 
* **The Systemic Failure:** Because CrowdStrike is security software that runs at the kernel level, the operating system could not isolate the error. This escalated a single fault into a massive failure, leaving 8.5 million independent computers stuck on a Blue Screen of Death.