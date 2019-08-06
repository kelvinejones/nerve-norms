# Nerve Norms

[Nerve Norms](https://nervenorms.bellstone.ca) is the website for the Nerve Excitability Test (NET). It is a **clinical tool** which was designed as part of a MSc thesis to **allow health care workers to measure the health of peripheral nerves**.

## Machine Learning and Analysis

A separate repository ([healthy-nerves](https://github.com/stellentus/healthy-nerves)) contains the code used for the analysis that led to this site. It included the following:
* **Imputing missing data** in health measurements (using linear regression, multiple imputation, and an autoencoder). Some of the nerve measurements are missing in some participants, so those values had to be filled before further analysis could be undertaken.
* **Measuring ethnic differences and batch effects in international data** to (with unsupervised clustering and variation of information). Previous studies suspected differences due to ethnicity, so international data could not be combined until those purported differences could be demonstrated or disproven.
* **Distinguishing unhealthy nerves from healthy for personalized medicine** using probabilistic modelling and outlier detection. The website compares each individual patient to data collected from hundreds of other healthy people to measure whether or not the patient's nerve is healthy.

## Project Layout

* The **JavaScript** (using **D3** for charts) in `res/templates/` is statically hosted at https://nervenorms.bellstone.ca.
* The **go** functions (`convert`, `norms`, `outliers`, and `participants`) are hosted as **Google Cloud Functions**.
* A few sample data files are in `res/data`, but the main data files are in a different repository.
