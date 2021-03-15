# Feature Spec - Portieris Operator 
**Is your feature request related to a problem? Please describe.**
- I want to have Portieris admission controller with Operator. 
- I want to install Portieris operator from Operator Hub.


**Describe the solution you'd like**
- Implement Go operator for Portieris admission controller. 
- Implement bundle and CSV for installation from OLM / Operator Hub.


**Describe alternatives you've considered**  
Helm Operator can be used to generate operator from Helm chart. 
However, the post-install annotation used in Portieris Helm chart for configuring MutatingAdmissionWebhook does not work in Helm Operator. Actually, Helm Operator reconcile logic reverts the post-install changes with object expected from dry-run, which does not consider post-install steps.


**Additional context**
- implemented Go-version operator with the latest Portieris code.
- confirmed it works on ROKS (OCP4.5) and OCP4.6/4.7.