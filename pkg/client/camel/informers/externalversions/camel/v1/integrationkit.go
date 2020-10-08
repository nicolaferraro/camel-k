/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	camelv1 "github.com/apache/camel-k/pkg/apis/camel/v1"
	versioned "github.com/apache/camel-k/pkg/client/camel/clientset/versioned"
	internalinterfaces "github.com/apache/camel-k/pkg/client/camel/informers/externalversions/internalinterfaces"
	v1 "github.com/apache/camel-k/pkg/client/camel/listers/camel/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// IntegrationKitInformer provides access to a shared informer and lister for
// IntegrationKits.
type IntegrationKitInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.IntegrationKitLister
}

type integrationKitInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewIntegrationKitInformer constructs a new informer for IntegrationKit type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewIntegrationKitInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredIntegrationKitInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredIntegrationKitInformer constructs a new informer for IntegrationKit type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredIntegrationKitInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CamelV1().IntegrationKits(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CamelV1().IntegrationKits(namespace).Watch(context.TODO(), options)
			},
		},
		&camelv1.IntegrationKit{},
		resyncPeriod,
		indexers,
	)
}

func (f *integrationKitInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredIntegrationKitInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *integrationKitInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&camelv1.IntegrationKit{}, f.defaultInformer)
}

func (f *integrationKitInformer) Lister() v1.IntegrationKitLister {
	return v1.NewIntegrationKitLister(f.Informer().GetIndexer())
}
