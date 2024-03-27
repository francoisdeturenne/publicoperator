# publicoperator
a podinfo operator for training sessions - helm repo and docker image as well

error msge in k8s according to chart:

6.5.1 --> 7.3.0 (pas dans l'index)  test-tp-deploy-sgsn-operator-hl
False                  HelmChart
'orange-cicd-ofr/orange-cicd-ofr-test-tp-deploy-sgsn-operator-hl' is not ready:
latest generation of object has not been reconciled   

6.5.1 --> 7.5.0 test-tp-deploy2-podinfo-hl                                           Unknown
Running 'upgrade' action with timeout of 5m0s
25m        
et 
au bout de 5 mn: test-tp-deploy2-podinfo-hl                      False
Helm upgrade failed for release test-tp-deploy2-podinfo/test-tp-deploy2-podinfo
with chart podinfo@7.5.0: context deadline exceeded                      29m

NB: NVEAUX pods PENDING 

7.5.0 -> 7.3.0 
test-tp-deploy2-podinfo-hl                                           Unknown
Running 'upgrade' action with timeout of 5m0s
32m   
et
au bout de 5 mn::kusto
test-tp-deploy2-podinfo-hl                      False                      Helm
upgrade failed for release test-tp-deploy2-podinfo/test-tp-deploy2-podinfo with
chart podinfo@7.3.0: context deadline exceeded        

NB: NVEAUX pods PENDING 

7.3.0 -> 6.5.1 echoue si on ne nettoie pas (helm uninstall / msge time out

6.5.1 -> 7.3.0 pending/timeout, replicaset pending 2/4 - helm list: 6.5.1
deployed/ 7.3.0 failed

7.3.0 -> 6.5.1 est-tp-deploy2-podinfo-hl                           True
Helm upgrade succeeded for release
test-tp-deploy2-podinfo/test-tp-deploy2-podinfo.v3 with chart podinfo@6.5.1 
version 7.3.0 uninstalled dans helm list -A replicas à 3 6.5.1 === 3/3

suppression deploy -> pods -> pas de drift detection ?


6.5.1 nettoyé -> 7.3.0: test-tp-deploy2-podinfo-hl                        False
Helm rollback to previous release
test-tp-deploy2-podinfo/test-tp-deploy2-podinfo.v7 with chart podinfo@6.5.1
succeeded               
apres
time out
retry
helm list -A = pas d'enreg deploy