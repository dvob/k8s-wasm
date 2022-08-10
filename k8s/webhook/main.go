package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gomodules.xyz/jsonpatch/v2"
	admit "k8s.io/api/admission/v1beta1"
	authn "k8s.io/api/authentication/v1beta1"
	authz "k8s.io/api/authorization/v1beta1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	http.HandleFunc("/authn", authenticationHandler)
	http.HandleFunc("/authz", authorizationHandler)
	http.HandleFunc("/admit", validatingAdmissionHandler)
	http.HandleFunc("/mut-admit", mutatingAdmissionHandler)

	err := http.ListenAndServeTLS(":6666", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func authenticationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("AUTHN: failed to read body: %s", err)
		return
	}

	log.Printf("AUTHN: INPUT: '%s'", body)

	tr := &authn.TokenReview{}
	err = json.Unmarshal(body, tr)
	if err != nil {
		log.Print("AUTHN:", err)
		return
	}

	resp := &authn.TokenReview{
		Status: authn.TokenReviewStatus{
			Authenticated: false,
			User:          authn.UserInfo{},
			Error:         "token is not magic",
		},
	}

	if tr.Spec.Token == "magic-token" {
		log.Print("AUTHN: token correct")
		resp.Status = authn.TokenReviewStatus{
			Authenticated: true,
			User: authn.UserInfo{
				UID:      "0",
				Username: "magic-user",
				Groups: []string{
					//"system:masters",
					"magic-group",
				},
			},
		}
	} else {
		log.Print("AUTHN: invalid token")
	}

	output, err := json.Marshal(resp)
	if err != nil {
		log.Print("AUTHN:", err)
		return
	}
	log.Printf("AUTHN: OUTPUT: '%s'", output)
	w.Write(output)
}

// authorizationHandler permits the group magic-group to edit configmaps with the name foo
func authorizationHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("AUTHZ: failed to read body: %s", err)
		return
	}

	log.Printf("AUTHZ: INPUT: '%s'", body)

	sar := &authz.SubjectAccessReview{}
	err = json.Unmarshal(body, sar)
	if err != nil {
		log.Print("AUTHZ:", err)
		return
	}

	resp := &authz.SubjectAccessReview{
		Status: authz.SubjectAccessReviewStatus{
			Allowed: false,
			Denied:  false,
		},
	}

	if sar.Spec.ResourceAttributes != nil &&
		sar.Spec.ResourceAttributes.Resource == "configmaps" &&
		contains(sar.Spec.Groups, "magic-group") {
		log.Print("AUTHZ: magic-group authorized for configmap")
		resp.Status.Allowed = true
	} else {
		log.Print("AUTHZ: not authorized")
	}

	output, err := json.Marshal(resp)
	if err != nil {
		log.Print("AUTHZ:", err)
		return
	}
	log.Printf("AUTHZ: OUTPUT: '%s'", output)
	w.Write(output)
}

func validatingAdmissionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("VALIDATE: failed to read body: %s", err)
		return
	}
	log.Printf("VALIDATE: INPUT: '%s'", body)

	ar := &admit.AdmissionReview{}
	err = json.Unmarshal(body, ar)
	if err != nil {
		log.Print("VALIDATE:", err)
		return
	}

	cm := &corev1.ConfigMap{}
	err = json.Unmarshal(ar.Request.Object.Raw, cm)
	if err != nil {
		log.Print("VALIDATE:", err)
		return
	}

	resp := &admit.AdmissionReview{
		Response: &admit.AdmissionResponse{
			UID:     ar.Request.UID,
			Allowed: true,
		},
	}
	if _, ok := cm.Data["not-allowed-value"]; ok {
		log.Print("VALIDATE: reject configmap with value not-allowed-value")
		resp.Response.Allowed = false
		resp.Response.Result = &v1.Status{
			Message: "value 'not-allowed-value' not allowed in configmap",
		}
	} else {
		log.Print("VALIDATE: configmap validated")
	}

	output, err := json.Marshal(resp)
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("VALIDATE: OUTPUT: '%s'", output)
	w.Write(output)

}

func mutatingAdmissionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("MUT: failed to read body: %s", err)
		return
	}

	log.Printf("MUT INPUT: '%s'", body)

	ar := &admit.AdmissionReview{}
	err = json.Unmarshal(body, ar)
	if err != nil {
		log.Print("MUT:", err)
		return
	}

	cm := &corev1.ConfigMap{}
	err = json.Unmarshal(ar.Request.Object.Raw, cm)
	if err != nil {
		log.Print("MUT:", err)
		return
	}

	if cm.Data == nil {
		cm.Data = map[string]string{}
	}
	cm.Data["magic-value"] = "foobar"

	target, err := json.Marshal(cm)
	if err != nil {
		log.Print(err)
		return
	}

	patches, err := jsonpatch.CreatePatch(ar.Request.Object.Raw, target)
	if err != nil {
		log.Print("MUT:", err)
		return
	}

	patch, err := json.Marshal(patches)
	if err != nil {
		log.Print("MUT:", err)
		return
	}

	log.Print("MUT: patch:", string(patch))

	patchType := admit.PatchTypeJSONPatch
	resp := &admit.AdmissionReview{
		Response: &admit.AdmissionResponse{
			UID:       ar.Request.UID,
			Allowed:   true,
			Patch:     patch,
			PatchType: &patchType,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Print(err)
		return
	}

	log.Printf("MUT: OUTPUT: '%s'", string(data))
	w.Write(data)
}

func contains(names []string, lookup string) bool {
	for _, name := range names {
		if name == lookup {
			return true
		}
	}
	return false
}
