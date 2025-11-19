## **Grounding Probability and Causality**

## 1. Introduction

The formal apparatus of causal modeling, as developed by Spirtes, Glymour, and Scheines, Pearl, and others, has provided a powerful framework for causal discovery and counterfactual prediction. This framework is built upon a precise, axiomatic connection between causality, usually represented as directed acyclic graphs, and joint probability distributions. The power of the framework lies in its ability to discover a class of causal structures from conditional independencies and to make predictions about the joint probability distribution under interventions on those structures. 

The entire inferential edifice of causal modeling stands upon three main principles: the Causal Markov (CM) condition, the Causal Minimality (CMin) condition, and the Faithfulness condition. In the standard interventionist framework, these conditions are defended as a bundle of foundational assumptions--justified as methodologically useful or intuitively plausible. While powerful, this leaves them without a coherent grounding.

The issue with missing grounds upon which causal modeling rests is compounded by an older problem: the interpretation of probabilities. The Kolmogorov axioms provide a formal calculus, but their interpretation remains a matter of dispute. The standard (non-subjectivist) interpretations of probability are frequentist, logical, and propensity theories. Frequentists identify evidence for probabilities with the probabilities themselves, which is a mistake; so, this leaves logical and propensity theories.

This paper argues that Abstract Object Theory (AOT), as detailed in the _Principia Logico-Metaphysica_ (PLM), provides a toolkit for grounding both probability and causality. At first glance, AOT invites a logical interpretation of probability. AOT provides a theory of possible worlds, defined as abstract objects that encode only propositional properties (e.g., _being such that_ $p$) and that possibly are such that all and only propositions they encode are true. Given this, one might attempt to define the probability of a proposition $p$ as a ratio of cardinalities of classes of possible worlds. This approach, however, is a non-starter in AOT since possible worlds are not discernible in the manner required for extensional counting.

Grounding both probability and causality in AOT, then, requires something like a propensity approach to probability, which understands probability to be a relation involving a generating setup, or, as we shall see, a situation. This paper attempts to show that this is feasible, and then, building upon this, it will provide grounding for _intervention_, _direct cause_, and _causal mechanism_. 

The central thesis of this paper is one of reductive unification. We argue that by grounding concepts of _probability_, _intervention_, _direct cause_, and _causal mechanism_ in AOT, CM and CMin emerge not as independent axioms, but as theorems that follow from these concepts. The project is not to prove that the physical universe is Markovian, but to demonstrate the internal logical coherence of the interventionist framework itself, showing that CM and CMin are entailments of its core commitments. Faithfulness, by contrast, will be shown to be a non-derivable, contingent assumption.

## 2. Formal Preliminaries

Abstract Object Theory (AOT) is a typed modal logic distinguished by its novel theory of predication (or two ways in which things can be said to bear relations). We start by summarizing only the essential components.

#### 2.1. Two Modes of Predication

The core of AOT is a fundamental distinction between two ways in which things can be said to bear relations. 

1. **Exemplification** Ordinary objects (e.g., a table) and abstract objects (e.g., the number two) exemplify properties. This is represented by standard predication $Fx$. For example, the number two is abstract.
2. **Encoding** Abstract objects also encode properties. This is represented by a new mode of predication, $xF$. For example, the number two is prime.

The key distinction is that an abstract object's identity is determined by the properties it encodes, not by the properties it exemplifies. This allows AOT to model abstract objects as entities that are constituted by a set of properties, without those properties needing to be true _of_ the object in the standard exemplification sense.

#### 2.2. The Comprehension Principle for Abstract Objects

AOT guarantees the existence of abstract objects via a powerful Comprehension Principle. For any condition $\phi$ on properties (where $x$ does not occur free), there exists an abstract object ($A!$) that encodes all and only the properties that satisfy $\phi$:

$$
\exists x (A!x \land \forall F(xF \equiv \phi))
$$

This axiom is the engine of AOT's ontology. It allows us to reify any description of a collection of properties into a single abstract object whose entire nature is given by that description. In addition, given that identity for abstract objects is given in terms of the properties they encode (rather than exemplify), it is provable in the system that the above axiom leads to a stronger claim of uniqueness. Finally, the above axiom generalizes to any type of abstract object from individuals of type $i$ to propositions of type $\langle\rangle$ and relations of type $\langle i, i\rangle$ and beyond.

#### 2.3. Situations

The machinery of AOT can be used to define _situations_. A situation $x$ is an abstract object that encodes only propositional properties (or properties like _being such that_ $p$ or $[\lambda y p]$, for some proposition $p$):

$$
\text{Situation}(x) \equiv_{df} A!x \land \forall F(xF \rightarrow \exists p(F = [\lambda y p]))
$$

Using the Comprehension Principle for Abstract Objects, it is easy to establish that (unique) situations exist. Finally, we say that a proposition $p$ is _true in_ $x$ ($x \models p$) just in case $x$ is a situation that encodes propositional property $[\lambda y p]$. This, along with the resulting theory of situations in AOT, provides the formal tools we need to ground probabilities and causality.

## 3. Propensity Theory of Probability

The literature on propensities, from Popper onwards, has been hampered by the vagueness of its core concepts, particularly the generating conditions that are said to possess propensity. The first constructive task is to build a propensity theory of probability. 

#### 3.1. Propensity

A generating condition or chance setup $c$ is a physical situation. For example, a fair die being tossed on a table is one such physical situation. In AOT, we can represent this setup not as a concrete particular, but as _the_ abstract situation that fully describes it. Using the Comprehension Principle for Abstract Objects, we can say, for example, that $s_c$ is _the_ situation where all and only propositions true of $c$ (i.e., $q_c$) are true in $s_c$. This $s_c$ is _the_ abstract object that is the generating conditions, reified.

Given that we will use abstract situations instead of generating conditions or chance setups, and given that we may understand the outcome of interest as a proposition, we now posit a primitive propends relation $P!$ that holds amongst a proposition and two abstract individuals.

**Propends Relation:** $P!$ is the 3-place propends relation. $P!s, p, r$ is being a situation, a proposition, and a rational number such that $s$ propends $p$ with magnitude $r$.

(Note that in the future, we could expand this to incorporate time, so that $P!$ would be the 5-place relation of being a situation $s$, a moment in time $t_1$, a proposition $p$, a moment in time $t_2$, and a rational number $r$ such that $s$ at $t_1$ propends $p$ at $t_2$ with magnitude $r$. This would probably make some of the arguments in the causality section more appealing, but, for now, we continue with the simpler case.)

#### 3.2. Grounding the Kolmogorov Axioms

The primitive propensity relation is not the mathematical function that appears in the Kolmogorov axioms. To get from the primitive relation $P!$ to the mathematical function (which we will denote with " $\hat{p}$ ") that satisfies the Kolmogorov axioms (and thereby ground the latter in the former), we need to perform two steps: first, we need to relativize $P!$'s first argument, and, second, we claim that Kolmogorov's axioms govern $P!$ so relativized.

###### 3.2.1 Relativizing $P!$

The problem with $P!$, as it stands, is that when $p$ is entirely irrelevant to $s$, $P!s,p,r$ will be false for any $r$. To see that this is a problem, suppose that _the_ situation of interest describes the fair die being tossed on a table, $s_c$. If we now take $p$ = being such that the price of bread is $2, then, intuitively, $P!s_c,p,r$ is false no matter what $r$ is since $s_c$ does not propend $p$ to any magnitude. However, the probability that $p$ is true given that the die lands on 6 is perfectly well-defined. 

So, to get to the mathematical function $\hat{p}$, we need a situation that is relevant to every proposition. In AOT, it is a theorem that there is precisely one such object: the actual world $w_{\alpha}$. Because $w_{\alpha}$ encodes the propositions of the maximal and consistent situation, we may relativize $P!$ to $w_{\alpha}$ to get to $\hat{p}$. However, because it will be important later to also relativize the mathematical function to non-actual, possible worlds, we will say that there is a collection of such functions $`\hat{p}_{w}`$, relativized by $w$. (In practice, these probability functions will likely be relativized to situations that are proper parts of $w_{\alpha}$, but we stick with the maximal situations that are possible worlds for present purposes.)

So relative to $w$, $P!$ is a function because $\forall p \exists! r(P!w, p, r)$ and we identify $\hat{p}_w$ with this relativized $P!$. To make this seem more familiar, we can also define the application of this function.

**Probability Function Application:** $`\hat{p}_{w}(p) =_{df} \iota r(P!w,p,r)`$. _The value of_ $`\hat{p}_{w}`$ _for_ $p$, is _the_ object $r$ such that $P!w,p,r$.

###### 3.2.2 Axioms of Relativized $P!$

We now show how these functions $\hat{p}_w$ are governed by Kolmogorov's axioms. We posit the Non-negativity and Additivity axioms as constraints on our primitive $P!$ relation. The Normalization axiom, however, can be derived as a theorem by positing a more fundamental bridge principle that explicitly connects AOT's encoding-as-truth machinery with our propends primitive $P!$.

1. **Non-negativity (Posited):** $\forall p, r(P!w, p, r \rightarrow r \geq 0)$.
2. **Additivity (Posited):** For any disjoint propositions $p_1$, $p_2$, if $P!w, p_1, r_1$ and $P!w, p_2, r_2$, then $P!w, p_1 \lor p_2, r_1 + r_2$.
3. **Bridge Principle (Posited):** $w \models p \rightarrow P!w, p, 1$. This axiom links AOT to the probability calculus. It states that if a possible world $w$ makes a proposition true, then the propensity of $p$ relative to $w$ is 1.
4. **Normalization (Derived Theorem):** $P!w, \top, 1$, where $\top$ is the tautological proposition. Proof (Sketch): The tautological proposition $\top$ is a theorem in AOT; so, it is necessarily true: $\Box \top$. It is a theorem in AOT that $\Box p \equiv \forall w (w \models p)$. So, it follows that $\forall w (w \models \top)$. From our Bridge Principle, it follows that $P!w, \top, 1$.  

With this, we have successfully defined probability functions $\hat{p}_w$ and grounded them in our primitive propends $P!$ relation. In addition, we should extend this to conditional probability, as this will be important for what is to come.

**Conditional Probability Function Application:** 

$$
\hat{p}_w(p | q) =_{df} {\iota r(P!w, p \land q, r) \over \iota r(P!w, q, r)}
$$

provided $`\hat{p}_w(q) =_{df} \iota r(P!w, q, r)  > 0`$.

#### 3.3. Humphreys's Paradox

Humphreys (1985) argues that propensities cannot be probabilities because propensities are asymmetric and probabilities are not. In our language, it makes no sense to say instead that $p$ propends $s$ with magnitude $r$ (i.e., $P!p, s, r$). On the other hand, in the probability calculus, if $\hat{p}_w(p | q)$ is defined, then so is $\hat{p}_w(q | p)$.

The above AOT framework simultaneously validates Humphreys's core insight while rejecting his conclusion. The paradox arises from a conflation of two distinct concepts, which our framework separates:

1. **Level 1:** The primitive $P!s, p, r$ is asymmetric. It is not causal, but might be seen as capturing causal facts. 
2. **Level 2:** The defined function $\hat{p}_w$, world-bound to $w$, is the standard, symmetric probability function obeying Kolmogorov's axioms.

There is no contradiction. $\hat{p}_w(p | q)$ is defined by a single value $r$, which is the ratio of _the_ object $r_1$ such that $P!w, p \land q, r_1$ and _the_ object $r_2$ such that $P!w, q, r_2$. $\hat{p}_w(q | p)$ is also well defined given the above discussion. AOT's ability to formally distinguish the underlying primitive relation $P!$ from the mathematical functions $\hat{p}_w$ resolves the paradox. 

## 4. Interventionist Theory of Causality

With a propensity theory of probability established, we now build upon it to ground an interventionist theory of causality. We shall formalize the central concepts of an interventionist theory--_intervention_, _direct cause_, and _causal mechanism_--within the current AOT framework. In what follows, it is convenient to take " $V = v$ " (the claim _that a variable takes a value_) to be another way of expressing a proposition.

**Intervention:** An intervention $I$ on $X$ with respect to $Y$, relative to ${\bf V}$, is an abstract relation (of type $\langle i, i\rangle$ guaranteed to exist by AOT's typed Comprehension Principle for abstract objects) that maps one possible world $w$ to another one $w'$, where there is a value $x$ of $X$ such that $w'$ makes true _being such that the probability of_ $X = x$ _is_ 1 _in_ $w'$, $w$ does not make true _being such that the probability of_ $X = x$ _is_ 1 _in_ $w$, and otherwise the two worlds are equivalent with respect to every proposition concerning ${\bf V} \backslash \{X, Y\}$:
1. $w' \models [\lambda z (\hat{p}_{w'}(X = x) = 1)]$ (This is equivalent to $w' \models X = x$, which by the Bridge Principle in section 3 implies $P!w',X = x, 1$).
2. $w \not\models [\lambda z (\hat{p}_{w}(X = x) = 1)]$.
3. $w \models p \equiv w' \models p$, with respect to every proposition concerning ${\bf V} \backslash \{X, Y\}$.

Given this, we can define "direct cause".

**Direct Cause:** A variable $X$ is a direct cause of variable $Y$ $(X \in \text{DC}(Y))$, relative to ${\bf V}$, if and only if there exists at least two distinct values $x_1, x_2$ for $X$ and some constant setting for all other variables $Z_1, ..., Z_n$ in ${\bf V} \backslash \{X, Y\}$ such that the probability of $Y$ differs between two corresponding worlds. That is, suppose we have world $w'$ resulting from an intervention on $X, Z_1, ..., Z_n$ with respect to $Y$ and world $w''$ resulting from an intervention on $X, Z_1, ..., Z_n$ with respect to $Y$ such that $w'$ and $w''$ differ only in the propositions they make true about $X$ (i.e., $X = x_1$ vs. $X = x_2$) with respect to ${\bf V} \backslash \{Y\}$. Then $X$ is a direct cause of $Y$ just in case $`\hat{p}_{w'}(Y = y) \neq \hat{p}_{w''}(Y = y)`$.

In the causal modeling literature, a mechanism is the process that determines the value (or probability distribution) of a variable from the value of its direct causes. In our AOT framework, this mechanism is the fact that a specific stable propends relation holds. 

**Causal Mechanism:** Let $X_1, ..., X_n$ be the direct causes of $Y$. Then 

$$
p_{\text{mech}_Y} =_{df} [\lambda z \forall x_1, ..., x_n(\hat{p}_w(Y = y | X_1 = x_1, ..., X_n = x_n) = r)].
$$ 

The causal mechanism for $Y$ is _being such that the probability of_ $Y = y$ _given_ $X_1 = x_1, ..., X_n = x_n$ _is_ $r$.

This definition grounds causal mechanism fully in AOT, but before moving on, we must make one more move. A core idea from the interventionist framework is modularity.

**Modularity (Posited):** Causal mechanisms are modular in the sense that an intervention $I$ on $X$ with respect to $Y$, relative to ${\bf V}$, does not alter the causal mechanisms $`p_{\text{mech}_Z}`$ for any other variable $Z \in {\bf V} \backslash \{X\}$. That is, $`w' \models p_{\text{mech}_Z} \equiv w \models p_{\text{mech}_Z}`$, where $w'$ is the resulting intervention world.

Note that this posited axiom governing causal mechanisms does not follow from the definition of "intervention". There, we said the two worlds do not differ with respect to propositions regarding ${\bf V}\backslash \{X, Y\}$. But with the modularity axiom, we have the stronger claim that the two worlds also do not differ on $p_{\text{mech}_Y}$.

## 5. Causal Markov, Causal Minimality, and Faithfulness

This section demonstrates that CM and CMin are theorems that follow from the definitions in the previous two sections. We also clarify the status of Faithfulness.

#### 5.1. Deriving CM

CM states that, relative to ${\bf V}$, every variable $X$ is probabilistically independent of its non-effects $\text{NE}(X) = {\bf V} \backslash \text{DC}(X)$ given its direct causes $\text{DC}(X)$. That is, the observational probability distribution is such that: 

$$
\hat{p}_{w}(X | \text{DC}(X), \text{NE}(X)) = \hat{p}_w(X | \text{DC}(X))
$$

**Proof (Analytic Entailment):** The theorem follows directly from our AOT-definition of "causal mechanism".

1. The causal mechanism $p_{\text{mech}_X}$ is the proposition that the probability of $X$ in $w$ is only a function of $\text{DC}(X)$.

2. The propensities generated by $w$ for any $X = x$ are therefore, by definition, sensitive only to the values of $\text{DC}(X)$.

3. The variables in $\text{NE}(X)$ do not include $\text{DC}(X)$.

4. Therefore, propositions about the values of $\text{NE}(X)$ are not part of $p_{\text{mech}_X}$.

5. Because the propensities for $X$ are determined only by $p_{\text{mech}_X}$, adding information about $\text{NE}(X)$ to the conditioning set is irrelevant.

6. Therefore, $\hat{p}_w(X | \text{DC}(X), \text{NE}(X)) = \hat{p}_w(X | \text{DC}(X))$.

#### 5.2. Deriving CMin

CMin states that for any direct cause $X$ of $Y$, $X$ is not a redundant cause of $Y$. $X$ is a redundant cause of $Y$ if $Y$ is observationally independent of $X$ conditional on all of $Y$'s other direct causes. 

**Proof (Sketch):** The proof of this theorem demonstrates that our definitions are logically incompatible with a failure of CMin.

1. Assume for contradiction that CMin is false.

2. This means there exists at least one direct cause $X$ of $Y$ that is redundant.

3. Let $Z_1, ..., Z_n$ be the other direct causes of $Y$.

4. The redundancy from step 2 implies that $Y$ is observationally independent of $X$ conditional on $Z_1, ..., Z_n$.

5. This implies that $\hat{p}_w(Y|X, Z_1, ..., Z_n) = \hat{p}_w(Y | Z_1, ..., Z_n)$.

6. This means that for any value $y$, any distinct values $x_1$, $x_2$ and any values $z_1, ..., z_n$, $\hat{p}_w(Y = y | X = x_1, Z_1 = z_1, ..., Z_n = z_n) = \hat{p}_w(Y = y | X = x_2, Z_1 = z_1, ..., Z_n = z_n)$.

7. However, by assumption 2, $X$ is a direct cause of $Y$. So, given the definition of "direct cause", there exists $x_1$, $x_2$ and some values $z_1, ..., z_n$ such that $`\hat{p}_{w'}(Y = y) \neq \hat{p}_{w''}(Y = y)`$, where $w'$ is the resulting world from an intervention $X = x_1, Z_1 = z_1, ..., Z_n = z_n$ and $w''$ is the resulting world from an intervention $X = x_2, Z_1 = z_1, ..., Z_n = z_n$.

8. Now we must bridge the interventional probabilities in step 7 with the observational probabilities in step 6.
    
    a. By our modularity axiom, the intervention world $w'$ is such that the causal mechanism for $Y$ still holds.

    b. By our definition of "causal mechanism", then, $\hat{p}_{w'}(Y | X, Z_1, ..., Z_n) = \hat{p}_w(Y|X, Z_1, ..., Z_n)$.

    c. In the intervention world $w'$, the values of $X$ and $Z_1, ..., Z_n$ are fixed. Therefore, the unconditional probability of $Y$ in $w'$ is, by definition, its probability conditional on those fixed facts: $`\hat{p}_{w'}(Y = y) = \hat{p}_{w'}(Y = y | X = x_1, Z_1 = z_1, ..., Z_n = z_n)`$.

    d. Substituting the result of 8c back into 8b gives us $`\hat{p}_{w'}(Y = y) = \hat{p}_{w}(Y = y | X = x_1, Z_1 = z_1, ..., Z_n = z_n)`$ and, similarly for $w''$, we get $`\hat{p}_{w''}(Y = y) = \hat{p}_{w}(Y = y | X = x_2, Z_1 = z_1, ..., Z_n = z_n)`$.

9. Substituting the result of step 8 back into the probabilities in step 7 gives us a contradiction with step 6.

10. Therefore, CMin is true.

#### 5.3. The Status of Faithfulness

This AOT framework also clarifies the status of the Faithfulness condition. Faithfulness states that the _only_ conditional independencies in the probability distribution are the ones entailed by CM. Unlike CM and CMin, Faithfulness is not a derivable theorem of AOT. Suppose $w_{\alpha}$ encodes a set of causal mechanism propositions that _coincidentally cancel each other out_. For example, we have not said anything that prevents $w_{\alpha}$ from encoding a mechanism for $Y$ such that $X$ and $Z$ are direct causes of $Y$, $X$ is also a direct cause of $Z$, but the two 'causal paths' perfectly cancel each other out, making $X$ probabilistically independent of $Y$ and thereby violating Faithfulness.

Faithfulness is therefore not a theorem in AOT, but rather an assumption about the nature of the specific situation being investigated. It is often a crucial bridge assumption required to move from an observed probability distribution to the underlying causal mechanisms, but our AOT framework makes its status as a non-derivable assumption explicit.

## 6. Conclusion

This paper constructs a unified foundation for both probability and causality, originating from AOT. 

The first constructive thesis was the development of a two-level theory of probability. We posited a _primitive asymmetric propends_ relation $P!$ that holds between situations that represent generating conditions, propositions as outcomes, and a rational number. This relation was then used to define a collection of standard probability functions $\hat{p}_w$. This two-level structure was shown to provide a resolution to Humphreys's Paradox by separating the asymmetric propends relation from the symmetric mathematics of probability.

The second constructive thesis built directly upon this propensity theory of probability to formalize the interventionist account of causation. We provided a definition of "intervention", "direct cause", and "causal mechanism" via AOT. Given this, we demonstrated the _entailment_ of foundational principles of causal modeling.

It is crucial to be precise about this contribution. To reiterate, the paper does not claim to have proven CM in a way that would satisfy critics who question CM's application to physical systems. Instead, we have shown that CM and CMin are _consequences_ of definitions and deeper assumptions. The paper's contribution is therefore a reductive unification: it replaces a loose bundle of assumptions with a single foundation that demonstrates the internal logical coherence of the interventionist framework such that CM and CMin are derivable. At the same time, it clarifies that Faithfulness is a separate, non-derivable assumption.
