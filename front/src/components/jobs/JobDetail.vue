<template>
  <b-card no-body>
    <b-tabs card>
      <b-tab title="pipeline" active>
        <b-container class="pipeline-container" fluid="true">
          <b-row align-h="center" align-v="center">
            <b-col>
              <b-card 
                bg-variant="success"
                class="mb-2 card job-card"
              >
              <b-container slot="header" class="step-header">
                <b-row align-h="end" align-v="center">
                  <b-col cols="4">
                    <h4 class="step-header-title">step 1</h4>
                  </b-col>
                  <b-col cols="4" align-self="end" class="step-actions">
                    <i class="fa fa-terminal"></i>
                    <i class="fa fa-cog"></i>
                  </b-col>
                </b-row>
              </b-container>
              <p class="card-text">
              fetch data from repo
              </p>
              </b-card>
            </b-col>
            <b-col md="auto">
              <i class="fa fa-arrow-right"></i>
            </b-col>
            <b-col>
              <b-card 
                bg-variant="success"
                class="mb-2 card job-card"
                >
                <b-container slot="header" class="step-header">
                  <b-row align-h="end" align-v="center">
                    <b-col cols="4">
                      <h4 class="step-header-title">step 2</h4>
                    </b-col>
                    <b-col cols="4" align-self="end" class="step-actions">
                      <i class="fa fa-terminal"></i>
                      <i class="fa fa-cog"></i>
                    </b-col>
                  </b-row>
                </b-container>
                <p class="card-text">
                build back code
                </p>
              </b-card>
            </b-col>
            <b-col md="auto">
              <i class="fa fa-arrow-right"></i>
            </b-col>
            <b-col>
              <b-card 
                bg-variant="success"
                class="mb-2 card job-card"
                >
                <b-container slot="header" class="step-header">
                  <b-row align-h="end" align-v="center">
                    <b-col cols="4">
                      <h4 class="step-header-title">step 3</h4>
                    </b-col>
                    <b-col cols="4" align-self="end" class="step-actions">
                      <i class="fa fa-terminal"></i>
                      <i class="fa fa-cog"></i>
                    </b-col>
                  </b-row>
                </b-container>
                <p class="card-text">
                test back code
                </p>
              </b-card>
            </b-col>
            <b-col md="auto">
              <i class="fa fa-arrow-right"></i>
            </b-col>
            <b-col>
              <b-card 
                bg-variant="success"
                class="mb-2 card job-card"
                >
                <b-container slot="header" class="step-header">
                  <b-row align-h="end" align-v="center">
                    <b-col cols="4">
                      <h4 class="step-header-title">step 4</h4>
                    </b-col>
                    <b-col cols="4" align-self="end" class="step-actions">
                      <i class="fa fa-terminal"></i>
                      <i class="fa fa-cog"></i>
                    </b-col>
                  </b-row>
                </b-container>
                <p class="card-text">
                  deploy app
                </p>
              </b-card>
            </b-col>
          </b-row>
        </b-container>
      </b-tab>
      <b-tab title="repository configuration" >
      </b-tab>
    </b-tabs>
  </b-card>
</template>
<script>
// todo co;pute cols props regarding of step number
</script>
<style scoped>
.pipeline-container {
  margin-top: 1rem;
}
.step-header-title {
  margin-bottom: 0;
}
.step-actions {
  display: flex;
  justify-content: end;
}
.step-actions .fa {
  font-size: 1.5rem;
  margin-right: 1rem;
  color: var(--darkest-color);
}
.step-actions .fa:last-of-type {
  margin-right: 0;
}
.step-actions .fa:hover {
  color: var(--dark-grey);
  cursor: pointer;
}
.fa-arrow-right {
  font-size: 5rem;
}
</style>
